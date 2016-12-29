angular.module('master', [])
    .controller('MasterController', ['$scope', '$http', '$state', '$cookieStore', 'Auth',
        function($scope, $http, $state, $cookieStore, Auth) {
            $scope.$watch(Auth.isLoggedIn, function(value, oldValue) {
                if (!value && oldValue) {
                    $state.go('login')
                }

                if (value) {
                    $http.get('/user', config)
                        .success(function(data, status) {
                            $scope.user = data
                        })
                        .error(function(data, status) {
                            console.log("something went logging in")
                            $state.go('login')
                        })
                }

            }, true);
            $scope.deployments = []
            $http.get('/deployments', config)
                .success(function(data, status) {
                    $scope.deployments = data
                })
                .error(function(data, status) {
                    console.log("something went wrong getting deployments")
                })
            $scope.logout = function() {
                $cookieStore.remove("auth")
            };
            $http.get('/user', config)
                .success(function(data, status) {
                    $scope.user = data
                })
                .error(function(data, status) {
                    console.log("something went logging in")
                    $state.go('login')
                })

        }
    ])