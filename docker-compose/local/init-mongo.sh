mongo -- "flyover" <<EOF
    var rootUser = 'root';
    var rootPassword = 'root';
    var admin = db.getSiblingDB('admin');
    admin.auth(rootUser, rootPassword);

    var user = 'root';
    var passwd = '$(cat "root")';
    db.createUser({user: user, pwd: passwd, roles: ["readWrite"]});
EOF