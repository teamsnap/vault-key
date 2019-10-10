{
    "targets": [
        {
            "target_name": "vault-auto",
            "conditions":[
                ["OS=='linux'", {
                    "sources": [ "vault-auto.linux.cc" ],
                    "libraries": [ "../vault.linux.a" ]
                }],
                ["OS=='mac'", {
                    "sources": [ "vault-auto.darwin.cc" ],
                    "libraries": [ "../vault.darwin.a" ]
                }]
            ],
        },
    ],
}
