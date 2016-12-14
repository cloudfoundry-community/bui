angular.module('deployments',[])
.controller('DeploymentsController', ['$scope','$http',
function ($scope,$http) {
  $http.get('/deployments', config)
  .success(function(data, status) {
    $scope.deployments = data
  })
  .error(function(data, status) {
     console.log("something went wrong getting deployments")
  })
}])
