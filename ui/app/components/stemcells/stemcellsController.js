angular.module('stemcells',[])
.controller('StemcellsController', ['$scope', '$http',
function ($scope, $http) {
  $http.get('/stemcells', config)
  .success(function(data, status) {
    $scope.stemcells = data
  })
  .error(function(data, status) {
     console.log("something went wrong getting stemcells")
  })
}])
