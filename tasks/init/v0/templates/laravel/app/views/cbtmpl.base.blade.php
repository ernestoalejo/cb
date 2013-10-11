<!DOCTYPE html>
<!--[if lt IE 7]>      <html class="no-js lt-ie9 lt-ie8 lt-ie7"> <![endif]-->
<!--[if IE 7]>         <html class="no-js lt-ie9 lt-ie8"> <![endif]-->
<!--[if IE 8]>         <html class="no-js lt-ie9"> <![endif]-->
<!--[if gt IE 8]><!--> <html class="no-js"> <!--<![endif]-->

<head>

  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">

  <title>@yield('title')</title>

  <!-- concat:css /styles/{{% .AppName %}}.css -->
    <link rel="stylesheet" href="/styles/bootstrap.css">
    <link rel="stylesheet" href="/styles/app.css">
  <!-- endconcat -->
  
  <style type="text/css"> .ng-cloak { display: none; } </style>

  @if(App::environment() !== 'local' && Config::get('analytics.google') !== '')
    <script>
      var _gaq=[['_setAccount','{{ Config::get('analytics.google') }}']];
      (function(d,t){var g=d.createElement(t),s=d.getElementsByTagName(t)[0];
      g.src=('https:'==location.protocol?'//ssl':'//www')+'.google-analytics.com/ga.js';
      s.parentNode.insertBefore(g,s)}(document,'script'));

      @if (!$admin)
        window._gaq.push(['_trackPageview']);
      @endif
    </script>
  @endif

  @yield('header')

</head>
<body ng-app="app" ng-controller="AppCtrl">

  <div class="navbar navbar-default navbar-fixed-top"><div class="container">
    <div class="navbar-header">
      <a href="#/" class="navbar-brand">{{% .AppName %}}</a>
    </div>
    <ul class="nav navbar-nav">
      <li class="active"><a href="#/">Home</a></li>
      <li><a href="#/other-path">Other Page</a></li>
    </ul>
  </div></div>

  <div class="container">

    <div ng-view></div>

  </div>

  <div id="loading" ng-show="GlobalMsg.isLoading()" style="display: none;"
      ng-animate="{show: 'loading-enter', hide: 'loading-leave'}">
    <div><img src="/images/spinner.gif" width="32" height="32"></div>
  </div>
  <div id="notifications" class="notifications"></div>

  <!-- concat:js /scripts/ie.js -->
    <!--[if lt IE 9]>
      <script src="/components/es5-shim/es5-shim.min.js"></script>
      <script src="/components/json3/lib/json3.min.js"></script>
      <script> var concat_script_here; </script>
    <![endif]-->
  <!-- endconcat -->

  <!-- min --><script src="//ajax.googleapis.com/ajax/libs/jquery/1.8.3/jquery.js"></script>
  <script>
    window.jQuery || document.write('<script src="/components/jquery-1.8.3/jquery.min.js">' +
      '<\/script>');
  </script>

  <!-- concat:js /scripts/{{% .AppName %}}.js -->
    <!-- min --><script src="/components/bower-angular/angular.js"></script>
    <!-- min --><script src="/components/bower-angular/angular-sanitize.js"></script>
    <!-- min --><script src="/components/jquery-placeholder/jquery.placeholder.js"></script>

    <!-- compile /scripts/vendor.js -->
      <script src="/components/angular-bootstrap/ui-bootstrap-tpls.js"></script>
      <script src="/components/bootstrap-notify/js/bootstrap-notify.js"></script>
      <script src="/components/bower-angular/i18n/angular-locale_es.js"></script>

      {{-- This ones should keep the order --}}
      <script src="/components/bootstrap/js/transition.js"></script>
      <script src="/components/bootstrap/js/alert.js"></script>
    <!-- endcompile -->

    <!-- compile /scripts/scripts.js -->
      <script src="/scripts/app.js"></script>
      <script src="/scripts/controllers/global.js"></script>
      <script src="/scripts/controllers/home.js"></script>
      <script src="/scripts/directives/match.js"></script>
      <script src="/scripts/directives/placeholder.js"></script>
      <script src="/scripts/error-handler.js"></script>
      <script src="/scripts/http-interceptor.js"></script>
      <script src="/scripts/services/global.js"></script>
      <script src="/scripts/services/location.js"></script>
      <script src="/scripts/services/modals.js"></script>
      <script src="/scripts/services/scroll.js"></script>
      <script src="/scripts/services/user.js"></script>
      <script src="/scripts/services/user.js"></script>
    <!-- endcompile -->
  <!-- endconcat -->

  @if (false)
    <!-- compile /scripts/test.js -->
      <script src="components/bower-angular/angular-mocks.js"></script>      
      <script src="scripts/app.test.js"></script>
    <!-- endcompile -->
  @endif

  <script type="text/javascript">
    @foreach($data as $d)
      angular.module('{{ $d['module'] }}').constant('{{ $d['name'] }}', {{ json_encode($d['value']) }});
    @endforeach
  </script>

  @yield('scripts')

</body>
</html>
