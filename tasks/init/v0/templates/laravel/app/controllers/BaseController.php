<?php

class BaseController extends Controller {

  /**
   * Setup the layout used by the controller.
   *
   * @return void
   */
  protected function setupLayout() {
    if (!is_null($this->layout))
    {
      $this->layout = View::make($this->layout);
    }
  }

  protected function buildBase() {
    // User info
    $admin = false;
    $name = '';
    $id = '';
    $user = Auth::user();
    if ($user) {
      $admin = $user->admin;
      $name = $user->username;
      $id = $user->id;
    }

    return View::make('base')->
      with('admin', true)->
      with('user', $user)->
      with('data', array(
        array(
          'module' => 'services.admin.user',
          'name' => 'admin',
          'value' => $admin,
        ),
        array(
          'module' => 'services.admin.user',
          'name' => 'name',
          'value' => $name,
        ),
        array(
          'module' => 'services.admin.user',
          'name' => 'id',
          'value' => $id,
        ),
      ));
  }

  protected function json($data) {
    $headers = array(
      'Cache-Control' => 'max-age=0,no-cache,no-store,post-check=0,pre-check=0',
      'Expires' => 'Mon, 26 Jul 1997 05:00:00 GMT',
      'Content-Type' => 'application/json; charset=utf-8',
      'X-Response-Processor' => 'json',
    );
    $content = ")]}',\n" . json_encode($data);
    return Response::make($content, 200, $headers);
  }

  protected function unsafeJson($data) {
    $headers = array(
      'Cache-Control' => 'max-age=0,no-cache,no-store,post-check=0,pre-check=0',
      'Expires' => 'Mon, 26 Jul 1997 05:00:00 GMT',
      'Content-Type' => 'application/json; charset=utf-8',
      'X-Response-Processor' => 'json',
    );
    $content = json_encode($data);
    return Response::make($content, 200, $headers);
  }

}
