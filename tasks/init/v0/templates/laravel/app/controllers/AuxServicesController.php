<?php

class AuxServicesController extends BaseController {

  public function notFound() {
    $data = \Validators\AuxServices\NotFound::validateJson();
    Log::error('client route not found: ' . $data['path']);
    return $this->json(array('success' => true));
  }

  public function reporter() {
    $data = \Validators\AuxServices\Reporter::validateJson();
    $name = $data['name'];
    $message = $data['message'];
    $stack = $data['stack'];
    Log::error("CLIENT ERROR:\n * Name: $name\n * Message: $message\n" .
        " * Stack: $stack\n\n");
    return $this->json(array('success' => true));
  }

}
