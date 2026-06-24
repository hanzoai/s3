package shell

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

type Case struct {
	args   []string
	result string
}

var (
	TestCases = []*Case{
		//add circuit breaker config for global
		{
			args: strings.Split("-global -type Count -actions Read,Write -values 500,200", " "),
			result: `{
				"global": {
					"enabled": true,
					"actions": {
						"Read:Count": 500,
						"Write:Count": 200
					}
				}
			}`,
		},

		//disable global config
		{
			args: strings.Split("-global -disable", " "),
			result: `{
				"global": {
					"actions": {
						"Read:Count": 500,
						"Write:Count": 200
					}
				}
			}`,
		},

		//add circuit breaker config for buckets x,y,z
		{
			args: strings.Split("-buckets x,y,z -type Count -actions Read,Write -values 200,100", " "),
			result: `{
				"global": {
					"actions": {
						"Read:Count": 500,
						"Write:Count": 200
					}
				},
				"buckets": {
					"x": {
						"enabled": true,
						"actions": {
							"Read:Count": 200,
							"Write:Count": 100
						}
					},
					"y": {
						"enabled": true,
						"actions": {
							"Read:Count": 200,
							"Write:Count": 100
						}
					},
					"z": {
						"enabled": true,
						"actions": {
							"Read:Count": 200,
							"Write:Count": 100
						}
					}
				}
			}`,
		},

		//disable circuit breaker config of x
		{
			args: strings.Split("-buckets x -disable", " "),
			result: `{
			  "global": {
				"actions": {
				  "Read:Count": 500,
				  "Write:Count": 200
				}
			  },
			  "buckets": {
				"x": {
				  "actions": {
					"Read:Count": 200,
					"Write:Count": 100
				  }
				},
				"y": {
				  "enabled": true,
				  "actions": {
					"Read:Count": 200,
					"Write:Count": 100
				  }
				},
				"z": {
				  "enabled": true,
				  "actions": {
					"Read:Count": 200,
					"Write:Count": 100
				  }
				}
			  }
			}`,
		},

		//delete circuit breaker config of x
		{
			args: strings.Split("-buckets x -delete", " "),
			result: `{
			  "global": {
				"actions": {
				  "Read:Count": 500,
				  "Write:Count": 200
				}
			  },
			  "buckets": {
				"y": {
				  "enabled": true,
				  "actions": {
					"Read:Count": 200,
					"Write:Count": 100
				  }
				},
				"z": {
				  "enabled": true,
				  "actions": {
					"Read:Count": 200,
					"Write:Count": 100
				  }
				}
			  }
			}`,
		},

		//configure the circuit breaker for the size of the uploaded file for bucket x,y
		{
			args: strings.Split("-buckets x,y -type MB -actions Write -values 1024", " "),
			result: `{
			  "global": {
				"actions": {
				  "Read:Count": 500,
				  "Write:Count": 200
				}
			  },
			  "buckets": {
				"x": {
				  "enabled": true,
				  "actions": {
					"Write:MB": 1073741824
				  }
				},
				"y": {
				  "enabled": true,
				  "actions": {
					"Read:Count": 200,
					"Write:Count": 100,
					"Write:MB": 1073741824
				  }
				},
				"z": {
				  "enabled": true,
				  "actions": {
					"Read:Count": 200,
					"Write:Count": 100
				  }
				}
			  }
			}`,
		},

		//delete the circuit breaker configuration for the size of the uploaded file of bucket x,y
		{
			args: strings.Split("-buckets x,y -type MB -actions Write -delete", " "),
			result: `{
			  "global": {
				"actions": {
				  "Read:Count": 500,
				  "Write:Count": 200
				}
			  },
			  "buckets": {
				"x": {
				  "enabled": true
				},
				"y": {
				  "enabled": true,
				  "actions": {
					"Read:Count": 200,
					"Write:Count": 100
				  }
				},
				"z": {
				  "enabled": true,
				  "actions": {
					"Read:Count": 200,
					"Write:Count": 100
				  }
				}
			  }
			}`,
		},

		//enable global circuit breaker config (without -disable flag)
		{
			args: strings.Split("-global", " "),
			result: `{
			  "global": {
				"enabled": true,
				"actions": {
				  "Read:Count": 500,
				  "Write:Count": 200
				}
			  },
			  "buckets": {
				"x": {
				  "enabled": true
				},
				"y": {
				  "enabled": true,
				  "actions": {
					"Read:Count": 200,
					"Write:Count": 100
				  }
				},
				"z": {
				  "enabled": true,
				  "actions": {
					"Read:Count": 200,
					"Write:Count": 100
				  }
				}
			  }
			}`,
		},

		//clear all circuit breaker config
		{
			args: strings.Split("-delete", " "),
			result: `{
			
			}`,
		},
	}
)

func TestCircuitBreakerShell(t *testing.T) {
	// stored is the canonical ZAP-encoded payload the filer would hold; the
	// persistence seam (LoadConfig/SaveConfig) chains it across invocations,
	// exactly as production reads/writes it via the filer.
	var stored []byte
	cmd := &commandS3CircuitBreaker{}
	LoadConfig = func(commandEnv *CommandEnv, dir string, file string, buf *bytes.Buffer) error {
		_, err := buf.Write(stored)
		return err
	}
	SaveConfig = func(commandEnv *CommandEnv, dir string, file string, content []byte) error {
		stored = append([]byte(nil), content...)
		return nil
	}

	for _, tc := range TestCases {
		var writeBuf bytes.Buffer
		// -apply so the canonical payload is persisted through SaveConfig.
		err := cmd.Do(append(tc.args, "-apply"), nil, &writeBuf)
		if err != nil {
			t.Fatal(err)
		}

		actual := make(map[string]interface{})
		if err := json.Unmarshal(writeBuf.Bytes(), &actual); err != nil {
			t.Fatal(err)
		}

		expect := make(map[string]interface{})
		if err := json.Unmarshal([]byte(tc.result), &expect); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Fatalf("unexpected result for args %v:\n got: %v\nwant: %v", tc.args, actual, expect)
		}
	}
}
