## Introduction

Gc-simple-home (GCSH) implements a pretty REST server for the Comelit ecosystem (Simple Home, Vedo and VIP).
You can use the GCSH to enable the Voice Control (such as Amazon Alexa or Google Assistant) over the Comelit Simple Home (see the IFTTT section).

## How it works
GCSH receives a POST call and executes a GET from the serial bridge.

## Configuration
The configuration file (**conf/configuration.json**) includes all system parameters (such as Comelit Serial Bridge address, module addresses, custom commands, and so on).
The configuration can change at runtime (the system will re-load it when something changes).
Every env (an env can be a room, i.e., the kitchen) can be addressed with more names.
You can set many names by separating them using the **|** symbol (i.e., **kitchen|kitchen one|small kitchen**).
Below an example, inside the light section, for correctly set light's parameters for the kitchen.

`"kitchen|kitchen one|small kitchen":
    "id":"3",
    "type":"env"
}`
NOTE: the id should be fetched from the serial bridge. You can use the browser and look at the URL request.
You can create your custom env or command.

GCSH loads the configuration at runtime.

### Most important parameters
* network: includes the listening port, certificates and token (it can be blank, but it is recommended to set it).
* domotica: it includes the **targetAddress** (the serial bridge address) . It also includes all "env" such as light, shutter, alarm, scenario and so on. You can create your custom env.

### Configuration example
`{
  "network": {
    "http-port": 60001,
    "https-port": 60002,
    "https-key-file": "conf/host.key",
    "https-cert-file": "conf/host.crt",
    "token":""
  },
  "domotica": {
    "systemParameters": {
      "targetAddress": "192.168.188.15"
    },
    "entities": {
      "light": {
        "commands": {
          "on": "http://{targetAddress}/user/action.cgi?type=light&{type}1={id}&_=1548916457530",
          "off": "http://{targetAddress}/user/action.cgi?type=light&{type}0={id}&_=1548916457530"
        },
        "env": {
          "ingresso": {
            "id": "1",
            "type": "env"
          },
          "cucina": {
            "id": "3",
            "type": "env"
          },
          "lavanderia": {
            "id": "7",
            "type": "env"
          },
          "bagno|bagno 1": {
            "id": "12",
            "type": "env"
          },
          "camera 1|camera matrimoniale": {
            "id": "5",
            "type": "env"
          },
          "cabina|cabina armadio": {
            "id": "6",
            "type": "env"
          },
          "camera 2": {
            "id": "11",
            "type": "env"
          },
          "salotto": {
            "id": "1",
            "type": "num"
          },
          "studio": {
            "id": "12",
            "type": "env"
          },
          "camera 3": {
            "id": "14",
            "type": "env"
          },
          "bagno mansarda|bagno 2": {
            "id": "13",
            "type": "env"
          },
          "ripostiglio mansarda|ripostiglio due": {
            "id": "12",
            "type": "env"
          },
          "garage": {
            "id": "14",
            "type": "env"
          },
          "bagno 3|bagno garage": {
            "id": "16",
            "type": "env"
          }
        }
      },
      "shutter": {
        "commands": {
          "up": "http://{targetAddress}/user/action.cgi?type=shutter&{type}1={id}&_=1548925715292",
          "down": "http://{targetAddress}/user/action.cgi?type=shutter&{type}0={id}&_=1548925715292"
        },
        "env": {
          "salotto": {
            "id": "2",
            "type": "env"
          },
          "cucina": {
            "id": "3",
            "type": "env"
          }
        }
      },
      "clima": {
        "commands": {
          "on": "http://{targetAddress}/user/action.cgi?type=other&{type}1={id}&_=1548926656038",
          "off": "http://{targetAddress}/user/action.cgi?type=other&{type}0={id}&_=1548926656050"
        },
        "env": {
          "home": {
            "id": "0",
            "type": "num"
          }
        }
      },
      "alarm": {
        "commands": {
          "on": "http://{targetAddress}/user/action.cgi?vedo=1&tot={id}&_=1548938060394",
          "off": "http://{targetAddress}/user/action.cgi?vedo=1&dis={id}&_=1548938060479"
        },
        "env": {
          "garage": {
            "id": "2",
            "type": ""
          },
          "casa": {
            "id": "4",
            "type": ""
          },
          "mansarda": {
            "id": "0",
            "type": ""
          }
        }
      },
      "scenario": {
        "commands": {
          "on": "http://{targetAddress}/user/action.cgi?scenario={id}&_=1548939193034"
        },
        "env": {
          "luci": {
            "id": "0",
            "type": ""
          },
          "salotto": {
            "id": "4",
            "type": ""
          },
          "mansarda": {
            "id": "5",
            "type": ""
          }
        }
      }
    }
  }
}
`
## REST
The  GCSH currentry supports only the POST.
Allowed parameters are:
* entity: the entity name such as light, shutter, scenario. (It can be customizable)
* cmd: the command such as on, off. (It can be customizable).
* token: it authenticates the request. It can be left black (dangerous)

### Example
`{
   "target":"kitchen",
   "token":"eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0",
   "cmd":"on",
   "entity":"light"
}`

## IFTTT
If This Then That, also known as IFTTT, is a free web-based service to create chains of simple conditional statements, called applets.
An applet is triggered by changes that occur within other web services such as Google Assistant, Amazon Alexa, Gmail, Facebook, Telegram, Instagram, or Pinterest.

### How to create an IFTTT applet
* Go to IFTTT.com
* At the top right, click your username New Applet.
* Click this.
* Search for "Google Assistant."
* Click Google Assistant .
* Choose a trigger.
* Set the URL request (e.g. https://myddns.com:/60001)
* Set the **POST** method
* Set **application/json** as Content Type
* Set the body
`{ "target":"{{TextField}}", "token":"my_token","cmd":"on", "entity":"light" }`
* Complete the trigger fields: Enter up to three ways to say your phrase. ...
* Click Create trigger.

Note: you can change/add the entity name according to the configuration JSON file.


## Build
You can build the code by using the GO Lang compiler.

### Build Examples
#### Linux Alpine
`go build --ldflags '-w -linkmode external -extldflags "-static"' -o gcsh-alpine main.go`

#### X86_64
`go build -o gcsh-alpine main.go`

## Execution
Run gcsh. It will look at the configuration.json file.
`gcsh`

## Docker
Start the container and expose the port:
`docker run -d -i -t -p 60001:60001 gcammarata/gcsh`