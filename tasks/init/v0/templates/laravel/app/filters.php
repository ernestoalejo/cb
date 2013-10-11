<?php

/*
|--------------------------------------------------------------------------
| Application & Route Filters
|--------------------------------------------------------------------------
|
| Below you will find the "before" and "after" events for the application
| which may be used to do any work before or after a request into your
| application. Here you may also register your custom route filters.
|
*/

App::before(function($request) {
  Session::regenerate();

  // Save & change the XSRF token
  $oldToken = Session::get('xsrf', '');
  $newToken = Str::random(32);
  Session::put('xsrf', $newToken);

  // We can't use Cookie::put because it sets HttpOnly cookies only
  setcookie('XSRF-TOKEN', Crypt::encrypt($newToken), 0, '/');

  if (Request::getMethod() !== 'GET' && !App::environment() === 'local') {
    // Check the cookie, the header, and the code to see if everything is
    // looking correct in the request
    $cookie = isset($_COOKIE['XSRF-TOKEN']) ? $_COOKIE['XSRF-TOKEN'] : '';
    $header = Request::header('x-xsrf-token');
    $header = ($header === null) ? '' : $header[0];
    if ($cookie != $header) {
      throw new Illuminate\Session\TokenMismatchException;
    }
    if (strlen($cookie) == 0) {
      throw new Illuminate\Session\TokenMismatchException;
    }

    $value = Crypt::decrypt($cookie);
    if ($value != $oldToken) {
      throw new Illuminate\Session\TokenMismatchException;
    }
  }
});


App::after(function($request, $response) {
  if ($response instanceof Response) {
    $response->header('X-UA-Compatible', 'chrome=1');
  }
});

/*
|--------------------------------------------------------------------------
| Authentication Filters
|--------------------------------------------------------------------------
|
| The following filters are used to verify that the user of the current
| session is logged into this application. The "basic" filter easily
| integrates HTTP Basic authentication for quick, simple checking.
|
*/

Route::filter('admin', function() {
  if (!Auth::user()->admin) {
    App::abort(403);
  }
});


Route::filter('auth', function() {
  if (Auth::guest()) {
    return Redirect::to('accounts/login');
  }
});

/*
Route::filter('auth.basic', function()
{
	return Auth::basic();
});*/

/*
|--------------------------------------------------------------------------
| Guest Filter
|--------------------------------------------------------------------------
|
| The "guest" filter is the counterpart of the authentication filters as
| it simply checks that the current user is not logged in. A redirect
| response will be issued if they are, which you may freely change.
|
*/
/*
Route::filter('guest', function()
{
	if (Auth::check()) return Redirect::to('/');
});
*/
/*
|--------------------------------------------------------------------------
| CSRF Protection Filter
|--------------------------------------------------------------------------
|
| The CSRF filter is responsible for protecting your application against
| cross-site request forgery attacks. If this special token in a user
| session does not match the one given in this request, we'll bail.
|
*/
/*
Route::filter('csrf', function()
{
	if (Session::token() != Input::get('_token'))
	{
		throw new Illuminate\Session\TokenMismatchException;
	}
});*/
