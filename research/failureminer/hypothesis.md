# Falsifiable Hypothesis

**Hypothesis**: Proactive SSD drain triggered by `FailureMiner` (using a 15% threshold on `Available Spare` depletion) will reduce VFS stall events by >40% compared to reactive replacement based on static `smartctl` error logs.

**Metrics**:
- **Recall**: Percentage of actual failures predicted before I/O timeout.
- **False-Alarm Rate**: Percentage of "predicted" failures where the drive remained healthy for >30 days.
- **Lead Time**: Average time between prediction and drive failure.
