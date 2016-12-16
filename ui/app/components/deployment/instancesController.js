angular.module('deploymentInstances', [])
    .controller('DeploymentInstancesController', ['$scope', '$http', '$stateParams',
        function($scope, $http, $stateParams) {

            $http.get('/deployments/' + $stateParams.name + "/vms", config)
                .success(function(data, status) {
                    $scope.instances = data
                })
                .error(function(data, status) {
                    console.log("something went wrong getting deployment")
                })


            $scope.getLabel = function(state) {
                if (state == "running") {
                    return "success"
                }
                return "danger"
            }

        }
    ])