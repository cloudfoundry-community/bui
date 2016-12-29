angular.module('dashboard', [])
    .controller('DashboardController', ['$scope', '$http','$timeout',
        function($scope, $http, $timeout) {
            $scope.runningTasks = function() {
                $http.get('/tasks/running', config)
                    .success(function(data, status) {
                        $scope.running_tasks = data
                    })
                    .error(function(data, status) {
                        console.log("something went wrong getting running tasks")
                    })
            }
            $http.get('/info', config)
                .success(function(data, status) {
                    $scope.info = data
                })
                .error(function(data, status) {
                    console.log("something went wrong getting info")
                })

            $scope.runningTasks();
            $scope.intervalFunction = function() {
                $timeout(function() {
                    $scope.runningTasks();
                    $scope.intervalFunction();
                }, 5000);
            };
            $scope.intervalFunction();
        }
    ])