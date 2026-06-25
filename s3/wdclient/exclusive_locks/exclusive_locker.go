package exclusive_locks

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/wdclient"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
)

const (
	RenewInterval     = 4 * time.Second
	SafeRenewInterval = 3 * time.Second
	InitLockInterval  = 1 * time.Second
)

type ExclusiveLocker struct {
	token        int64
	lockTsNs     int64
	isLocked     atomic.Bool
	masterClient *wdclient.MasterClient
	lockName     string
	message      string
	clientName   string
	// Each lock has and only has one goroutine
	renewGoroutineRunning atomic.Bool
}

func NewExclusiveLocker(masterClient *wdclient.MasterClient, lockName string) *ExclusiveLocker {
	return &ExclusiveLocker{
		masterClient: masterClient,
		lockName:     lockName,
	}
}

func (l *ExclusiveLocker) IsLocked() bool {
	return l.isLocked.Load()
}

func (l *ExclusiveLocker) GetToken() (token int64, lockTsNs int64) {
	for time.Unix(0, atomic.LoadInt64(&l.lockTsNs)).Add(SafeRenewInterval).Before(time.Now()) {
		// wait until now is within the safe lock period, no immediate renewal to change the token
		time.Sleep(100 * time.Millisecond)
	}
	return atomic.LoadInt64(&l.token), atomic.LoadInt64(&l.lockTsNs)
}

func (l *ExclusiveLocker) RequestLock(clientName string) {
	if l.isLocked.Load() {
		return
	}

	// retry to get the lease
	for {
		if err := l.masterClient.WithZapClient(func(client *masterwire.Client) error {
			resp, err := client.LeaseAdminToken(context.Background(), masterwire.LeaseAdminTokenRequestInput{
				PreviousToken:    atomic.LoadInt64(&l.token),
				PreviousLockTime: atomic.LoadInt64(&l.lockTsNs),
				LockName:         l.lockName,
				ClientName:       clientName,
			})
			if err == nil {
				atomic.StoreInt64(&l.token, resp.Token())
				atomic.StoreInt64(&l.lockTsNs, resp.LockTsNs())
			}
			return err
		}); err != nil {
			glog.V(2).Infof("Failed to acquire lock %s: %v", l.lockName, err)
			time.Sleep(InitLockInterval)
		} else {
			break
		}
	}

	l.isLocked.Store(true)
	l.clientName = clientName
	glog.V(1).Infof("Acquired lock %s", l.lockName)

	// Each lock has and only has one goroutine
	if l.renewGoroutineRunning.CompareAndSwap(false, true) {
		// start a goroutine to renew the lease
		go func() {
			for {
				if l.isLocked.Load() {
					if err := l.masterClient.WithZapClient(func(client *masterwire.Client) error {
						resp, err := client.LeaseAdminToken(context.Background(), masterwire.LeaseAdminTokenRequestInput{
							PreviousToken:    atomic.LoadInt64(&l.token),
							PreviousLockTime: atomic.LoadInt64(&l.lockTsNs),
							LockName:         l.lockName,
							ClientName:       l.clientName,
							Message:          l.message,
						})
						if err == nil {
							atomic.StoreInt64(&l.token, resp.Token())
							atomic.StoreInt64(&l.lockTsNs, resp.LockTsNs())
							glog.V(2).Infof("Renewed lock %s: ts %d token %d", l.lockName, l.lockTsNs, l.token)
						}
						return err
					}); err != nil {
						glog.Warningf("Failed to renew lock %s: %v", l.lockName, err)
						l.isLocked.Store(false)
						return
					} else {
						time.Sleep(RenewInterval)
					}
				} else {
					time.Sleep(RenewInterval)
				}
			}
		}()
	}

}

func (l *ExclusiveLocker) ReleaseLock() {
	l.isLocked.Store(false)
	l.clientName = ""

	l.masterClient.WithZapClient(func(client *masterwire.Client) error {
		client.ReleaseAdminToken(context.Background(), masterwire.ReleaseAdminTokenRequestInput{
			PreviousToken:    atomic.LoadInt64(&l.token),
			PreviousLockTime: atomic.LoadInt64(&l.lockTsNs),
			LockName:         l.lockName,
		})
		return nil
	})
	atomic.StoreInt64(&l.token, 0)
	atomic.StoreInt64(&l.lockTsNs, 0)
}

func (l *ExclusiveLocker) SetMessage(message string) {
	l.message = message
}
