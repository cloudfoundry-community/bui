angular.module('bui', [
    'ui.ace',
    'ui.router',
    'ngCookies',
    'master',
    'dashboard',
    'auth',
    'login',
    'deployment',
    'deploymentManifest',
    'deploymentInstances',
    'deploymentSSH',
    'deploymentIndex',
    'deployments',
    'stemcells',
    'releases'
]).
config(['$stateProvider', '$urlRouterProvider',
    function config($stateProvider, $urlRouterProvider) {
        'use strict';
        $urlRouterProvider.otherwise('/login');
        $stateProvider
            .state('app', {
              abstract: true,
              templateUrl: 'app/components/dashboard/app.html'
            })
            .state('app.dashboard', {
                url: '/',
                templateUrl: 'app/components/dashboard/dashboard.html',
                controller: 'DashboardController'
            })
            .state('login', {
                url: '/login',
                templateUrl: 'app/components/login/login.html',
                controller: 'LoginController'
            })
            .state('app.deployments', {
                url: '/deployments/',
                templateUrl: 'app/components/deployments/deployments.html',
                controller: 'DeploymentsController'
            })
            .state('app.deployment', {
                url: '^/deployments/:name',
                templateUrl: 'app/components/deployment/deployment.html',
                controller: 'DeploymentController'
            })
            .state('app.deployment.index', {
                url: '/index',
                templateUrl: 'app/components/deployment/index.html',
                controller: 'DeploymentIndexController'
            })
            .state('app.deployment.manifest', {
                url: '/manifest',
                templateUrl: 'app/components/deployment/manifest.html',
                controller: 'DeploymentManifestController'
            })
            .state('app.deployment.instances', {
                url: '/instances',
                templateUrl: 'app/components/deployment/instances.html',
                controller: 'DeploymentInstancesController'
            })
            .state('app.deployment.ssh', {
                url: '/vms/:vm_name/ssh',
                templateUrl: 'app/components/deployment/ssh.html',
                controller: 'DeploymentSSHController'
            })
            .state('app.stemcells', {
                url: '/stemcells/',
                templateUrl: 'app/components/stemcells/stemcells.html',
                controller: 'StemcellsController'
            })
            .state('app.releases', {
                url: '/releases/',
                templateUrl: 'app/components/releases/releases.html',
                controller: 'ReleasesController'
            })
    }
]);