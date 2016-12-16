angular.module('master', [])
    .controller('MasterController', ['$scope', '$http', '$state','$cookieStore',
        function($scope, $http, $state, $cookieStore) {
            $scope.deployments = []
            $http.get('/deployments', config)
                .success(function(data, status) {
                    $scope.deployments = data
                })
                .error(function(data, status) {
                    console.log("something went wrong getting deployments")
                })
            $http.get('/user', config)
                .success(function(data, status) {
                    $scope.user = data
                })
                .error(function(data, status) {
                    $scope.logged_in = false
                    console.log("something went logging in")
                    $state.go('login')
                })
                $scope.logout = function() {
                  $cookieStore.remove("auth")
                };

        }
    ])