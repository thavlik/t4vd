import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';

import '../api.dart';
import '../model.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _usernameController = TextEditingController();
  final _passwordController = TextEditingController();
  final _usernameFocus = FocusNode();
  final _passwordFocus = FocusNode();
  bool _passwordVisible = false;
  bool _errorState = false;
  String? _usernameError;
  String? _passwordError;
  bool _isLoggingIn = false;

  bool get isUsernameValid => _usernameError == null;
  bool get isPasswordValid => _passwordError == null;

  void onUsernameChanged(String value) =>
      setState(() => _usernameError = usernameError(value));
  void onPasswordChanged(String value) =>
      setState(() => _passwordError = passwordError(value));

  String? usernameError(String value) {
    if (value.length < minUsernameLength || value.contains(" ")) {
      return 'Please enter your username';
    }
    return null;
  }

  String? passwordError(String value) {
    if (value.isEmpty) {
      return 'Please enter your password';
    }
    return null;
  }

  bool get isValid => isUsernameValid && isPasswordValid;

  void login(BuildContext context) async {
    if (_isLoggingIn) return;
    onUsernameChanged(_usernameController.text);
    onPasswordChanged(_passwordController.text);
    if (!isValid) {
      setState(() => _errorState = true);
      return;
    }
    setState(() {
      _errorState = false;
      _isLoggingIn = true;
    });
    final model = ScopedModel.of<BJJModel>(context);
    await model.login(
      username: _usernameController.text,
      password: _passwordController.text,
    );
    if (!mounted) return;
    setState(() => _isLoggingIn = false);
    Navigator.of(context).pop();
  }

  void signUp(BuildContext context) async {
    await Navigator.of(context).pushNamed("/user/signup");
    if (!mounted) return;
    if (!ScopedModel.of<BJJModel>(context).isLoggedIn) return;
    Navigator.of(context).pop();
  }

  void forgotPassword(BuildContext context) =>
      Navigator.of(context).pushNamed("/user/resetpassword");

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.all(12.0),
            child: Center(
              child: Container(
                constraints: const BoxConstraints(maxWidth: 300),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Container(
                      constraints: const BoxConstraints(maxWidth: 192),
                      child: const AspectRatio(
                        aspectRatio: 1.0,
                        child: Image(image: AssetImage('assets/laselva.png')),
                      ),
                    ),
                    const SizedBox(height: 32),
                    Padding(
                      padding: const EdgeInsets.all(8.0),
                      child: TextField(
                        key: const Key('signInUsername'),
                        enableSuggestions: false,
                        autocorrect: false,
                        controller: _usernameController,
                        focusNode: _usernameFocus,
                        textInputAction: TextInputAction.next,
                        autofocus: false,
                        onChanged: onUsernameChanged,
                        onSubmitted: (value) {
                          onUsernameChanged(value);
                          if (!isUsernameValid) {
                            _usernameFocus.requestFocus();
                            return;
                          }
                          _passwordFocus.requestFocus();
                        },
                        decoration: InputDecoration(
                          border: const OutlineInputBorder(),
                          hintText: 'Username',
                          labelText: _usernameError,
                          labelStyle: TextStyle(
                            color: _errorState ? Colors.red : Colors.white,
                          ),
                        ),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.all(8.0),
                      child: TextField(
                        key: const Key('signInPassword'),
                        obscureText: !_passwordVisible,
                        enableSuggestions: false,
                        autocorrect: false,
                        controller: _passwordController,
                        textInputAction: TextInputAction.done,
                        focusNode: _passwordFocus,
                        onSubmitted: (value) {
                          onPasswordChanged(value);
                          if (!isPasswordValid) {
                            _passwordFocus.requestFocus();
                            return;
                          }
                          login(context);
                        },
                        onChanged: onPasswordChanged,
                        decoration: InputDecoration(
                            border: const OutlineInputBorder(),
                            hintText: 'Password',
                            labelText: _passwordError,
                            labelStyle: TextStyle(
                              color: _errorState ? Colors.red : Colors.white,
                            ),
                            errorMaxLines: 3,
                            suffixIcon: IconButton(
                              onPressed: () {
                                setState(() {
                                  _passwordVisible = !_passwordVisible;
                                });
                              },
                              icon: Icon(_passwordVisible
                                  ? Icons.visibility
                                  : Icons.visibility_off),
                            )),
                      ),
                    ),
                    const SizedBox(height: 8),
                    ScopedModelDescendant<BJJModel>(
                        builder: (context, child, model) {
                      return SizedBox(
                        height: 24,
                        child: Visibility(
                          visible: model.loginErr != null,
                          child: Text(model.loginErr ?? '',
                              style: Theme.of(context)
                                  .textTheme
                                  .bodySmall!
                                  .copyWith(
                                    color: Colors.red,
                                  )),
                        ),
                      );
                    }),
                    const SizedBox(height: 8),
                    Opacity(
                      opacity: isValid && !_isLoggingIn ? 1.0 : 0.5,
                      child: Container(
                        decoration: BoxDecoration(
                          color: Theme.of(context).primaryColor,
                          borderRadius: BorderRadius.circular(2),
                        ),
                        child: Material(
                          color: Colors.transparent,
                          child: InkWell(
                            key: const Key('signInSubmit'),
                            onTap: () => login(context),
                            child: Row(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                Padding(
                                  padding: const EdgeInsets.all(16.0),
                                  child: Text(
                                    _isLoggingIn ? 'Logging in...' : 'Login',
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),
                    ),
                    const SizedBox(height: 20),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Tooltip(
                          message: 'Sign up',
                          child: TextButton(
                            key: const Key('signup'),
                            onPressed: () => signUp(context),
                            child: const Text('Sign up'),
                          ),
                        ),
                        TextButton(
                          onPressed: () => forgotPassword(context),
                          child: const Text('Forgot password?'),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
