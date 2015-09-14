function TaskCtrl($scope, $http) {
  $scope.tasks = [];
  $scope.working = false;

  var logError = function(data, status) {
    console.log('code '+status+': '+data);
    $scope.working = false;
  };

  var refresh = function() {
    return $http.get('/task/').
      success(function(data) { $scope.tasks = data.Tasks; }).
      error(logError);
  };

  $scope.addTask = function() {
    $scope.working = true;
    $http.post('/task/', {Title: $scope.todoText}).
      error(logError).
      success(function() {
        refresh().then(function() {
          $scope.working = false;
          $scope.todoText = '';
        })
      });
  };

  $scope.toggleCompleted = function(task) {
    data = {ID: task.ID, Title: task.Title, Completed: !task.Completed}
    $http.put('/task/'+task.ID, data).
      error(logError).
      success(function() { task.Completed = !task.Completed });
  };

  refresh().then(function() { $scope.working = false; });
}
