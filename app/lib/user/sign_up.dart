import 'dart:async';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:email_validator/email_validator.dart';
import 'package:scoped_model/scoped_model.dart';

import '../api.dart';
import '../model.dart';

class SignUpPage extends StatefulWidget {
  const SignUpPage({super.key});

  @override
  State<SignUpPage> createState() => _SignUpPageState();
}

class _SignUpPageState extends State<SignUpPage> {
  final _usernameController = TextEditingController();
  final _emailController = TextEditingController();
  final _firstNameController = TextEditingController();
  final _lastNameController = TextEditingController();
  final _passwordController = TextEditingController();
  final _passwordConfirmController = TextEditingController();

  final _usernameFocus = FocusNode();
  final _emailFocus = FocusNode();
  final _firstNameFocus = FocusNode();
  final _lastNameFocus = FocusNode();
  final _passwordFocus = FocusNode();
  final _passwordConfirmFocus = FocusNode();

  String? _usernameError;
  String? _emailError;
  String? _firstNameError;
  String? _lastNameError;
  String? _passwordError;
  String? _passwordConfirmError;

  bool _passwordVisible = false;
  Timer? _checkTakenUsername;
  Timer? _checkTakenEmail;
  final Map<String, bool> _takenUsernames = {};
  final Map<String, bool> _takenEmails = {};

  bool get _showAppBar => !kIsWeb;

  void submit() async {
    onUsernameChanged(_usernameController.text);
    onEmailChanged(_emailController.text);
    onFirstNameChanged(_firstNameController.text);
    onLastNameChanged(_lastNameController.text);
    onPasswordChanged(_passwordController.text);
    onPasswordConfirmChanged(_passwordConfirmController.text);
    if (!isValid) {
      return;
    }
    _checkTakenUsername?.cancel();
    _checkTakenEmail?.cancel();
    if (await usernameTaken() != false) {
      return;
    }
    if (await emailTaken() != false) {
      return;
    }
    if (!mounted) return;
    await ScopedModel.of<BJJModel>(context).register(
      username: _usernameController.text,
      email: _emailController.text,
      firstName: _firstNameController.text,
      lastName: _lastNameController.text,
      password: _passwordController.text,
    );
    if (!mounted) return;
    if (!ScopedModel.of<BJJModel>(context).isLoggedIn) return;
    Navigator.of(context).pop();
  }

  @override
  void dispose() {
    super.dispose();
    _checkTakenUsername?.cancel();
    _checkTakenEmail?.cancel();
  }

  Future<bool?> emailTaken() async {
    if (!mounted) {
      return true;
    }
    _checkTakenEmail?.cancel();
    final value = _emailController.text;
    bool? taken;
    try {
      taken = await emailExists(value);
      _takenEmails[value] = taken;
      setState(() => _emailError = emailError(value));
    } catch (err) {
      setState(() => _emailError = 'Error: ${err.toString()}');
    }
    return taken;
  }

  Future<bool?> usernameTaken() async {
    if (!mounted) {
      return true;
    }
    _checkTakenUsername?.cancel();
    final value = _usernameController.text;
    bool? taken;
    try {
      taken = await userExists(value);
      _takenUsernames[value] = taken;
      setState(() => _usernameError = usernameError(value));
    } catch (err) {
      setState(() => _usernameError = 'Error: ${err.toString()}');
    }
    return taken;
  }

  void onUsernameChanged(String value) {
    setState(() => _usernameError = usernameError(value));
    _checkTakenUsername?.cancel();
    if (value.length >= minUsernameLength &&
        !_takenUsernames.containsKey(value)) {
      _checkTakenUsername = Timer(const Duration(seconds: 1), usernameTaken);
    }
  }

  void onEmailChanged(String value) {
    setState(() {
      _emailError = emailError(value);
      if (_emailError == null) {
        _checkTakenEmail?.cancel();
        if (value.isNotEmpty && !_takenEmails.containsKey(value)) {
          _checkTakenEmail = Timer(const Duration(seconds: 1), emailTaken);
        }
      }
    });
  }

  void onFirstNameChanged(String value) =>
      setState(() => _firstNameError = firstNameError(value));
  void onLastNameChanged(String value) =>
      setState(() => _lastNameError = lastNameError(value));
  void onPasswordChanged(String value) =>
      setState(() => _passwordError = passwordError(value));
  void onPasswordConfirmChanged(String value) =>
      setState(() => _passwordConfirmError = passwordConfirmError(value));

  bool get isUsernameValid => _usernameError == null;
  bool get isEmailValid => _emailError == null;
  bool get isFirstNameValid => _firstNameError == null;
  bool get isLastNameValid => _lastNameError == null;
  bool get isPasswordValid => _passwordError == null;
  bool get isPasswordConfirmValid => _passwordConfirmError == null;

  bool get isValid =>
      isUsernameValid &&
      isEmailValid &&
      isFirstNameValid &&
      isLastNameValid &&
      isPasswordValid &&
      isPasswordConfirmValid;

  String? usernameError(String value) {
    if (value.length < minUsernameLength) {
      return 'Must contain at least $minUsernameLength characters';
    } else if (' \t\n!@#\$%^&*()\'".,?;:<>[]{}/\\'
        .split('')
        .any((c) => value.contains(c))) {
      return 'Cannot contain spaces or special characters';
    } else if (_takenUsernames[value] == true) {
      return 'Username already taken';
    }
    return null;
  }

  String? emailError(String value) {
    if (!EmailValidator.validate(value)) {
      return 'Enter a valid email';
    }
    return null;
  }

  String? firstNameError(String value) {
    if (value.isEmpty) {
      return 'Please enter your first name';
    }
    return null;
  }

  String? lastNameError(String value) {
    if (value.isEmpty) {
      return 'Please enter your last name';
    }
    return null;
  }

  String? passwordError(String value) {
    if (value.length < minPasswordLength) {
      return 'Must contain at least $minPasswordLength characters';
    }
    if (!value.contains(RegExp('[a-z]'))) {
      return 'Must contain lowercase letters';
    }
    if (!value.contains(RegExp('[A-Z]'))) {
      return 'Must contain uppercase letters';
    }
    if (!value.contains(RegExp('[0-9]'))) {
      return 'Must contain numbers';
    }
    if (!'!@#\$%^&*()\'".,?;:<>[]{}/\\'
        .split('')
        .any((c) => value.contains(c))) {
      return 'Must contain special characters';
    }
    return null;
  }

  String? passwordConfirmError(String value) {
    if (value != _passwordController.text) {
      return 'Passwords do not match';
    }
    return null;
  }

  void privacyPolicy(BuildContext context) =>
      Navigator.of(context).pushNamed('/user/privacypolicy');

  @override
  Widget build(BuildContext context) {
    final valid = isValid;
    return Scaffold(
      appBar: _showAppBar
          ? AppBar(
              title: const Text('Create an Account'),
              leading: IconButton(
                icon: const Icon(Icons.navigate_before),
                onPressed: () {
                  FocusScope.of(context).unfocus();
                  Navigator.of(context).pop();
                },
              ),
            )
          : null,
      body: Stack(
        children: [
          Center(
            child: SingleChildScrollView(
              child: Center(
                child: Container(
                  constraints: const BoxConstraints(maxWidth: 300),
                  child: SizedBox(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.center,
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        if (!_showAppBar)
                          Padding(
                            padding: const EdgeInsets.all(16.0),
                            child: Text(
                              "Create an Account",
                              style: Theme.of(context).textTheme.headlineSmall,
                            ),
                          ),
                        Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: TextField(
                            key: const Key('signUpUsername'),
                            enableSuggestions: false,
                            autocorrect: false,
                            controller: _usernameController,
                            autofocus: true,
                            focusNode: _usernameFocus,
                            textInputAction: TextInputAction.next,
                            onChanged: onUsernameChanged,
                            onSubmitted: (value) {
                              onUsernameChanged(value);
                              if (!isUsernameValid) {
                                _usernameFocus.requestFocus();
                              } else {
                                _emailFocus.requestFocus();
                              }
                            },
                            decoration: InputDecoration(
                              border: const OutlineInputBorder(),
                              hintText: 'Username',
                              labelText: () {
                                if (_usernameError != null) {
                                  return _usernameError;
                                }
                                if (_takenUsernames
                                    .containsKey(_usernameController.text)) {
                                  if (_takenUsernames[
                                          _usernameController.text] ==
                                      true) {
                                    return 'Username is already taken';
                                  } else {
                                    return 'Username is available';
                                  }
                                }
                                if (_usernameController.text.isEmpty) {
                                  return null;
                                }
                                return 'Checking availability';
                              }(),
                              labelStyle: TextStyle(
                                color: () {
                                  if (_takenUsernames
                                      .containsKey(_usernameController.text)) {
                                    if (_takenUsernames[
                                            _usernameController.text] ==
                                        true) {
                                      return Colors.red;
                                    } else {
                                      return Colors.green;
                                    }
                                  }
                                  return isUsernameValid
                                      ? Colors.white
                                      : Colors.red;
                                }(),
                              ),
                            ),
                          ),
                        ),
                        Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: TextField(
                            key: const Key('signUpEmail'),
                            enableSuggestions: false,
                            autocorrect: false,
                            controller: _emailController,
                            focusNode: _emailFocus,
                            textInputAction: TextInputAction.next,
                            onChanged: onEmailChanged,
                            onSubmitted: (value) {
                              onEmailChanged(value);
                              if (!isEmailValid) {
                                _emailFocus.requestFocus();
                              } else {
                                _firstNameFocus.requestFocus();
                              }
                            },
                            decoration: InputDecoration(
                              border: const OutlineInputBorder(),
                              hintText: 'Email',
                              labelText: () {
                                if (!isEmailValid) {
                                  return _emailError;
                                }
                                if (_takenEmails
                                    .containsKey(_emailController.text)) {
                                  if (_takenEmails[_emailController.text] ==
                                      true) {
                                    return 'Email is already taken';
                                  } else {
                                    return 'No other accounts with this email';
                                  }
                                }
                                if (_emailController.text.isEmpty) {
                                  return null;
                                }
                                return 'Checking availability';
                              }(),
                              labelStyle: TextStyle(
                                color: () {
                                  if (_takenEmails
                                      .containsKey(_emailController.text)) {
                                    if (_takenEmails[_emailController.text] ==
                                        true) {
                                      return Colors.red;
                                    } else {
                                      return Colors.green;
                                    }
                                  }
                                  return isEmailValid
                                      ? Colors.white
                                      : Colors.red;
                                }(),
                              ),
                            ),
                          ),
                        ),
                        Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: TextField(
                            key: const Key('signUpFirstName'),
                            enableSuggestions: false,
                            autocorrect: false,
                            controller: _firstNameController,
                            focusNode: _firstNameFocus,
                            textInputAction: TextInputAction.next,
                            onChanged: onFirstNameChanged,
                            onSubmitted: (value) {
                              onFirstNameChanged(value);
                              if (!isFirstNameValid) {
                                _firstNameFocus.requestFocus();
                              } else {
                                _lastNameFocus.requestFocus();
                              }
                            },
                            decoration: InputDecoration(
                              border: const OutlineInputBorder(),
                              hintText: 'First name',
                              labelText: _firstNameError,
                              labelStyle: const TextStyle(color: Colors.red),
                            ),
                          ),
                        ),
                        Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: TextField(
                            key: const Key('signUpLastName'),
                            enableSuggestions: false,
                            autocorrect: false,
                            controller: _lastNameController,
                            focusNode: _lastNameFocus,
                            textInputAction: TextInputAction.next,
                            onChanged: onLastNameChanged,
                            onSubmitted: (value) {
                              onLastNameChanged(value);
                              if (!isLastNameValid) {
                                _lastNameFocus.requestFocus();
                              } else {
                                _passwordFocus.requestFocus();
                              }
                            },
                            decoration: InputDecoration(
                              border: const OutlineInputBorder(),
                              labelText: _lastNameError,
                              labelStyle: const TextStyle(color: Colors.red),
                              hintText: 'Last name',
                            ),
                          ),
                        ),
                        Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: TextField(
                            key: const Key('signUpPassword'),
                            obscureText: !_passwordVisible,
                            enableSuggestions: false,
                            autocorrect: false,
                            controller: _passwordController,
                            textInputAction: TextInputAction.done,
                            focusNode: _passwordFocus,
                            onChanged: onPasswordChanged,
                            onSubmitted: (value) {
                              onPasswordChanged(value);
                              if (!isPasswordValid) {
                                _passwordFocus.requestFocus();
                              } else {
                                _passwordConfirmFocus.requestFocus();
                              }
                            },
                            decoration: InputDecoration(
                                border: const OutlineInputBorder(),
                                hintText: 'Password',
                                labelText: _passwordError,
                                labelStyle: const TextStyle(color: Colors.red),
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
                        Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: TextField(
                            key: const Key('signUpPasswordConfirm'),
                            obscureText: !_passwordVisible,
                            enableSuggestions: false,
                            autocorrect: false,
                            controller: _passwordConfirmController,
                            textInputAction: TextInputAction.done,
                            focusNode: _passwordConfirmFocus,
                            onChanged: onPasswordConfirmChanged,
                            onSubmitted: (value) {
                              onPasswordConfirmChanged(value);
                              if (!isPasswordConfirmValid) {
                                _passwordConfirmFocus.requestFocus();
                                return;
                              }
                              FocusScope.of(context).requestFocus(FocusNode());
                              submit();
                            },
                            decoration: InputDecoration(
                                border: const OutlineInputBorder(),
                                hintText: 'Confirm Password',
                                labelText: _passwordConfirmError,
                                labelStyle: const TextStyle(color: Colors.red),
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
                        const SizedBox(height: 20),
                        Tooltip(
                          message: 'Create account',
                          child: Opacity(
                            opacity: valid ? 1.0 : 0.5,
                            child: Container(
                              decoration: BoxDecoration(
                                color: Theme.of(context).primaryColor,
                                borderRadius: BorderRadius.circular(2),
                              ),
                              child: Material(
                                color: Colors.transparent,
                                child: InkWell(
                                  onTap: submit,
                                  child: Row(
                                    mainAxisAlignment: MainAxisAlignment.center,
                                    children: const [
                                      Padding(
                                        padding: EdgeInsets.all(16.0),
                                        child: Text("Sign Up"),
                                      ),
                                    ],
                                  ),
                                ),
                              ),
                            ),
                          ),
                        ),
                        const SizedBox(height: 20),
                        TextButton(
                          onPressed: () => privacyPolicy(context),
                          child: const Text("Privacy Policy"),
                        ),
                        const SizedBox(height: 64),
                      ],
                    ),
                  ),
                ),
              ),
            ),
          ),
          //Positioned(
          //  top: 16,
          //  left: 16,
          //  child: IconButton(
          //    icon: const Icon(Icons.navigate_before),
          //    onPressed: () => Navigator.of(context).pop(),
          //  ),
          //),
        ],
      ),
    );
  }
}
