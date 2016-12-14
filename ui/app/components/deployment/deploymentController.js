angular.module('deployment',[])
.controller('DeploymentController', ['$scope','$http', '$stateParams', '$state',
function ($scope, $http, $stateParams, $state) {
  $scope.name = $stateParams.name
}])
