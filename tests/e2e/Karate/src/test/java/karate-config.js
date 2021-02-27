function fn() {
    var env = karate.env; // get system property 'karate.env'
    karate.log('karate.env system property was:', env);
    if (!env) {
        env = 'dev';
    }
    var config = {
        env: env,
        baseUrl: 'http://localhost:8080/',
        adminAuthToken: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjEyMyJ9.eyJhdWQiOiJpdGxhYiIsImlzcyI6Imh0dHBzOi8vZGV2LmlkZW50aXR5LnJ0dWl0bGFiLnJ1IiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MDU0Njc3NTY4NjksInN1YiI6IjMyMS1sb2wtMzIxIiwiaXRsYWIiOlsicmVwb3J0cy5hZG1pbiIsInJlcG9ydHMudXNlciJdLCJzY29wZSI6WyJyb2xlcyIsIm9wZW5pZCIsInByb2ZpbGUiLCJpdGxhYi5ldmVudHMiLCJpdGxhYi5yZXBvcnRzIl19.N-rLG1WyhtUA-FwyDtebs4GkncvKMxcEUeQlN4elkc8',
        userAuthToken: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjEyMyJ9.eyJhdWQiOiJpdGxhYiIsImlzcyI6Imh0dHBzOi8vZGV2LmlkZW50aXR5LnJ0dWl0bGFiLnJ1IiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MDU0Njc3NTY4NjksInN1YiI6IjMyMS1sb2wtMzIxIiwiaXRsYWIiOiJ1c2VyIiwic2NvcGUiOlsicm9sZXMiLCJvcGVuaWQiLCJwcm9maWxlIiwiaXRsYWIuZXZlbnRzIiwiaXRsYWIucmVwb3J0cyJdfQ.YhVzNIoyfkqs_fJ8qHLLP7qJ1YWdsSUe4JgRTSff-6g',
        noITLabUserAuthToken: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjEyMyJ9.eyJhdWQiOiJpdGxhYiIsImlzcyI6Imh0dHBzOi8vZGV2LmlkZW50aXR5LnJ0dWl0bGFiLnJ1IiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MDU0Njc3NTY4NjksInN1YiI6IjMyMS1sb2wtMzIxIiwiaXRsYWIiOiJyZXBvcnRzLnVzZXIiLCJzY29wZSI6WyJyb2xlcyIsIm9wZW5pZCIsInByb2ZpbGUiLCJpdGxhYi5ldmVudHMiLCJpdGxhYi5yZXBvcnRzIl19.-kF-Zcxc6j3PYHLGHBqBnWIt1HIEZtoePt-waVqZ7zo'
    }
    if (env == 'dev') {
        // customize
    } else if (env == 'e2e') {
        // customize
        config.baseUrl = 'http://test-api:8080/';
    }
    return config;
}