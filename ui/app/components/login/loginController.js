angular.module('login',[])
.controller('LoginController', ['$scope','$http','$state','$window','Auth',
function ($scope, $http, $state, $window, Auth) {
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
        Auth.setUser(data)
        console.log("login")
        console.log(data)
        $state.go('app.dashboard')
      })
      .error(function(data, status) {
        console.log("fail")
        $window.location.href = '#/login?error=unauthorized'
      })

  }
}])
