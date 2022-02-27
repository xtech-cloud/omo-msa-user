package config

const defaultJson string = `{
	"service": {
		"address": ":7068",
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
		"name": "userCloud",
		"ip": "172.16.10.31",
		"port": "27017",
		"user": "root",
		"password": "pass2019",
		"type": "mongodb"
	}
}
`
