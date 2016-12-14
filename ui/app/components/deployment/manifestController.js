angular.module('deploymentManifest',[])
.controller('DeploymentManifestController', ['$scope','$http', '$stateParams',
function ($scope, $http, $stateParams) {
  console.log("MANIFEST BRO")
  $http.get('/deployments/' + $stateParams.name , config)
  .success(function(data, status) {
    console.log(data)
    $scope.manifest = data.manifest
  })
  .error(function(data, status) {
     console.log("something went wrong getting deployment")
  })
}])
