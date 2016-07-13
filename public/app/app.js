'use strict';

angular.module('ffhb', [
	'ngTable',
	'ngResource',
	'angularLocalStorage',
	'ui.router',
	'ui.bootstrap',
	'ui-leaflet',
	'angularMoment',
	'Authentication',
	'angular-web-notification',
	'config'
	])
	.config(['$urlRouterProvider',function ($urlRouterProvider){
		//,$httpProvider) {
		$urlRouterProvider.otherwise('/nodes/sort');
		//$locationProvider.html5Mode(true).hashPrefix('!');
		//$httpProvider.defaults.withCredentials = true;
	}]).run(function(amMoment,$cookieStore,$rootScope,$http) {
		amMoment.changeLocale('de');
		$rootScope.globals = $cookieStore.get('globals') || {};
		if ($rootScope.globals.currentUser) {
			$http.defaults.headers.common['Authorization'] = 'Basic ' + $rootScope.globals.currentUser.authdata; // jshint ignore:line
		}
	});