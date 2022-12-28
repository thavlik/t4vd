import 'package:flutter/material.dart';

class ResetPasswordPage extends StatefulWidget {
  const ResetPasswordPage({super.key});

  @override
  State<ResetPasswordPage> createState() => ResetPasswordPageState();
}

class ResetPasswordPageState extends State<ResetPasswordPage> {
  //final _usernameController = TextEditingController();
  //final _usernameFocus = FocusNode();
  //String? _usernameError;
  //bool get isUsernameValid => _usernameError == null;
  //void onUsernameChanged(String value) => setState(() => _usernameError = usernameError(value));

  String? usernameError(String value) {
    if (value.isEmpty) {
      return 'Please enter your username';
    }
    return null;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Reset Password'),
      ),
      body: Stack(
        children: [
          SingleChildScrollView(
            child: Padding(
              padding: const EdgeInsets.all(12.0),
              child: Center(
                child: Container(
                  constraints: const BoxConstraints(
                    maxHeight: 300,
                  ),
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Text(
                          "Password reset is not yet implemented. Please contact the system adminstrator to reset your password.",
                          style: Theme.of(context).textTheme.bodyLarge,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
