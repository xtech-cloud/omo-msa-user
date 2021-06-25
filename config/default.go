package config

const defaultJson string = `{
	"service": {
		"address": ":7078",
		"ttl": 15,
		"interval": 10
	},
	"logger": {
		"level": "info",
		"file": "logs/server.log",
		"std": false
	},
	"root":{
		"name":"root@admin",
		"psw":"aca44e364d2084c49fec383ea958eae2"
	},
	"database": {
		"name": "rgsCloud",
		"ip": "127.0.0.1",
		"port": "27017",
		"user": "root",
		"password": "pass2019",
		"type": "mongodb"
	}
}
`
