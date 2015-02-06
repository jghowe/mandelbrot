angular.module('mandelbrot', ['panzoom'])
  .controller('PanzoomController',['$scope', function($scope) {
    $scope.model = {};
    $scope.config = {};
  }]);
