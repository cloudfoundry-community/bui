angular.module('releases',[])
.controller('ReleasesController', ['$scope','$http',
function ($scope, $http) {
  $http.get('/releases', config)
  .success(function(data, status) {
    $scope.releases = data
  })
  .error(function(data, status) {
     console.log("something went wrong getting releases")
  })
  $scope.getLabel = function(state) {
    if (state) {
      return "success"
    }
    return "danger"
  }
}])
