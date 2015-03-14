angular.module('mandelbrot', ['panzoom', 'panzoomwidget'])
  .directive('mbLoad', ['$parse', function ($parse) {
    return {
      restrict: 'A',
      link: function (scope, elem, attrs) {
        var fn = $parse(attrs.mbLoad);
        elem.on('load', function (event) {
          scope.$apply(function() {
            fn(scope, { $event: event });
          });
        });
      }
    };
  }])
  .controller('PanzoomController',['$scope', 'PanZoomService', function($scope, PanZoomService) {
    $scope.model = {};
    $scope.config = {
        zoomLevels: 12,
        neutralZoomLevel: 5,
        scalePerZoomLevel: 1.5,
        zoomToFitZoomLevelFactor: 1
    };
    $scope.zoom = 1;
    $scope.centerX = -0.5;
    $scope.centerY = 0.0;
    $scope.width = 900;
    $scope.height = 600;

    $scope.getPosition = function() {
      PanZoomService.getAPI('panzoom').then(function (api) {
        var position = api.getModelPosition({x: $scope.width/2, y: $scope.height/2});
        var baseZoom = $scope.zoom;
        var zoom = Math.pow($scope.config.scalePerZoomLevel, $scope.model.zoomLevel - $scope.config.neutralZoomLevel);
        var width = zoom * $scope.width;
        var height = zoom * $scope.height;

        $scope.centerX = $scope.centerX + (((position.x - ($scope.width/2)) / $scope.width) / baseZoom * 3.0);
        $scope.centerY = $scope.centerY + (((position.y - ($scope.height/2)) / $scope.height) / baseZoom * 2.0);
        $scope.zoom = baseZoom * zoom;
      });
    }

    $scope.onImgLoad = function(event) {
      PanZoomService.getAPI('panzoom').then(function (api) {
        api.zoomToFit({x: 0, y: 0, width: $scope.width, height: $scope.height });
      });
    }
  }]);
