angular.module('deploymentIndex', [])
    .controller('DeploymentIndexController', ['$scope', '$http', '$stateParams',
        function($scope, $http, $stateParams) {

            $scope.instances_count = 0
            $http.get('/deployments/' + $stateParams.name, config)
                .success(function(data, status) {
                    $scope.manifest = YAML.parse(data.manifest)
                    for (var i = 0; i < $scope.manifest.jobs.length; i++) {
                        $scope.instances_count += $scope.manifest.jobs[i].instances;
                    }
                    console.log($scope.manifest)
                })
                .error(function(data, status) {
                    console.log("something went wrong getting deployment")
                })

                $scope.getLifecycle = function(lifecycle) {
                  if (lifecycle == "errand") {
                    return lifecycle
                  }
                  return "service"
                }
        }
    ])