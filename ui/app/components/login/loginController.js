angular.module('login',[])
.controller('LoginController', ['$scope','$http','$window',
function ($scope, $http, $window) {
  $scope.login = function(user) {
      var data = $.param({
          auth_type: "basic",
          username: user.name,
          password: user.password
      });
      var config = {
          headers: {
              'Content-Type': 'application/x-www-form-urlencoded'
          }
      }
      $http.post('/login', data, config)
      .success(function(data, status) {
        console.log("login")
        $window.location.href = '#/dashboard'
      })
      .error(function(data, status) {
        console.log("fail")
        $window.location.href = '#/login?error=unauthorized'
      })

  }
}])
