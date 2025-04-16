// src/pages/login.tsx

import { AuthPage } from "@refinedev/antd";
import { useLogin } from "@refinedev/core";

export const Login = () => {
  const { mutate: login, isLoading } = useLogin();

  return (
    <AuthPage
      type="login"
      formProps={{
        onFinish: async (values) => {
          try {
            const response = await login(values);

            // Explicitly cast to expected return shape
            const result = response as unknown as { success?: boolean; error?: any };

            if (result.success) {
              console.log("Login successful:", result);
            } else {
              console.error("Login failed:", result.error);
            }
          } catch (error) {
            console.error("Login exception:", error);
          }
        },
        initialValues: { email: "", password: "" },
      }}
    />
  );
};
