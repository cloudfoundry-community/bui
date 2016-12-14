angular.module('master',[])
.controller('MasterController', ['$scope', '$http',
function ($scope, $http) {
  $scope.deployments = []
  $http.get('/deployments', config)
  .success(function(data, status) {
    $scope.deployments = data
  })
  .error(function(data, status) {
     console.log("something went wrong getting deployments")
  })
}])
