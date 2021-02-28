function fn() {
    var env = karate.env; // get system property 'karate.env'
    karate.log('karate.env system property was:', env);
    if (!env) {
        env = 'dev';
    }

    var createJwt = function(subject, claims) {
        var JwtHelper = Java.type('JwtHelper');
        return JwtHelper.getAuthToken('itlab', subject, claims);
    }
    var createUserJwt = function(subject) {
        var JwtHelper = Java.type('JwtHelper');
        return JwtHelper.getAuthToken('itlab', subject, { itlab: ['user'] });
    }

    var config = {
        env: env,
        baseUrl: 'http://localhost:8080/',
        createJwt: createJwt,
        createUserJwt: createUserJwt,
        users: {
            admin: {
                id: 'user-admin',
                accessToken: createJwt('user-admin', { itlab: ['user', 'reports.admin'] })
            },
            plain: {
                id: 'user-plain',
                accessToken: createUserJwt('user-plain')
            },
            incorrect: {
                id: 'user-incorrect',
                accessToken: createJwt('user-admin', { itlab: ['reports.user'] })
            }
        }
    }
    if (env == 'dev') {
        // customize
    } else if (env == 'e2e') {
        // customize
        config.baseUrl = 'http://test-api:8080/';
    }
    return config;
}