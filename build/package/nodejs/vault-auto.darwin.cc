#include "vault.darwin.h"
#include <node.h>
#include <v8.h>
using namespace std;

namespace vaultc {

  using v8::FunctionCallbackInfo;
  using v8::Isolate;
  using v8::Local;
  using v8::MaybeLocal;
  using v8::NewStringType;
  using v8::Object;
  using v8::String;
  using v8::Value;
  using v8::Number;
  using v8::Exception;

  void secrets(const FunctionCallbackInfo<Value>& args) {
    Isolate* isolate = args.GetIsolate();

    // check number of args passed in
    if (args.Length() < 1) {
      // throw error that gets passed back to nodejs
      isolate->ThrowException(Exception::TypeError(
          String::NewFromUtf8(isolate, "Wrong number of arguments", v8::NewStringType::kNormal).ToLocalChecked()));
      return;
    }

    // check argument type
    if (!args[0]->IsString()) {
      isolate->ThrowException(Exception::TypeError(
          String::NewFromUtf8(isolate, "Wrong argument type, it should be a string", v8::NewStringType::kNormal).ToLocalChecked()));
      return;
    }

    v8::String::Utf8Value val(args[0]->ToString());
    std::string str (*val);

    // call golang function
    string secretData = GetSecrets((char*)str.c_str());

    args.GetReturnValue().Set(String::NewFromUtf8(
      isolate, secretData.c_str(), NewStringType::kNormal).ToLocalChecked());
  }

  void init(Local<Object> exports) {
    NODE_SET_METHOD(exports, "secrets", secrets);
  }

  NODE_MODULE(vault, init)
}
