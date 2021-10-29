package config

const defaultYAML string = `
service: 
    name: xtc.ogm.actor
    address: :18810
    ttl: 15
    interval: 10
logger:
    level: info
    dir: /var/log/ogm/
database:
    # 驱动类型，可选值为 [sqlite,mysql]
    driver: sqlite
    mysql:
        address: localhost:3306
        user: root
        password: mysql@XTC
        db: ogm
    sqlite:
        path: /tmp/ogm-actor.db
`
