function fn() {
    var env = karate.env; // get system property 'karate.env'
    karate.log('karate.env system property was:', env);
    if (!env) {
        env = 'dev';
    }
    var config = {
        env: env,
        baseUrl: 'http://localhost:8080/',
        users: {
            admin: {
                id: 'user-admin',
                accessToken: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjEyMyJ9.eyJhdWQiOiJpdGxhYiIsImlzcyI6Imh0dHBzOi8vZGV2LmlkZW50aXR5LnJ0dWl0bGFiLnJ1IiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MDU0Njc3NTY4NjksInN1YiI6InVzZXItYWRtaW4iLCJpdGxhYiI6WyJyZXBvcnRzLmFkbWluIiwidXNlciJdLCJzY29wZSI6WyJyb2xlcyIsIm9wZW5pZCIsInByb2ZpbGUiLCJpdGxhYi5ldmVudHMiLCJpdGxhYi5yZXBvcnRzIl19.Wksaz3faLzdEkWMaAjC0gkYT3SxRuHeLrlQordlQKes'
            },
            plain: {
                id: 'user-plain',
                accessToken: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjEyMyJ9.eyJhdWQiOiJpdGxhYiIsImlzcyI6Imh0dHBzOi8vZGV2LmlkZW50aXR5LnJ0dWl0bGFiLnJ1IiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MDU0Njc3NTY4NjksInN1YiI6InVzZXItcGxhaW4iLCJpdGxhYiI6InVzZXIiLCJzY29wZSI6WyJyb2xlcyIsIm9wZW5pZCIsInByb2ZpbGUiLCJpdGxhYi5ldmVudHMiLCJpdGxhYi5yZXBvcnRzIl19.hYPp4w6zJnL1ZQlO7TdaU7Jjvxz8hdF81uxjR0Xa6UY'
            },
            incorrect: {
                id: 'user-incorrect',
                accessToken: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjEyMyJ9.eyJhdWQiOiJpdGxhYiIsImlzcyI6Imh0dHBzOi8vZGV2LmlkZW50aXR5LnJ0dWl0bGFiLnJ1IiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MDU0Njc3NTY4NjksInN1YiI6InVzZXItaW5jb3JyZWN0IiwiaXRsYWIiOiJyZXBvcnRzLnVzZXIiLCJzY29wZSI6WyJyb2xlcyIsIm9wZW5pZCIsInByb2ZpbGUiLCJpdGxhYi5ldmVudHMiLCJpdGxhYi5yZXBvcnRzIl19.VNq05e-GsOX1QFOxQfom0uQq3xBLYzOHEf5wuDI7pM8'
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