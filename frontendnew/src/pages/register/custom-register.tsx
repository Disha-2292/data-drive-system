import React from "react";
import { Form, Input, Button, message, Typography } from "antd";

export const CustomRegister = () => {
  const [form] = Form.useForm();

  const onFinish = async (values: { name: string; email: string; password: string }) => {
    try {
      const response = await fetch("http://localhost:8080/auth/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(values),
      });

      const data = await response.json();

      if (response.ok) {
        message.success("Registration successful. Please login.");
        window.location.href = "/login";
      } else {
        message.error(data.error || "Registration failed");
      }
    } catch (error) {
      message.error("An unexpected error occurred");
    }
  };

  return (
    <div style={{ maxWidth: 400, margin: "50px auto" }}>
      <Typography.Title level={3}>Register</Typography.Title>
      <Form
        layout="vertical"
        form={form}
        onFinish={onFinish}
        initialValues={{ name: "", email: "", password: "" }}
      >
        <Form.Item name="name" label="Name" rules={[{ required: true }]}>
          <Input placeholder="Enter your name" />
        </Form.Item>
        <Form.Item name="email" label="Email" rules={[{ required: true, type: "email" }]}>
          <Input placeholder="Enter your email" />
        </Form.Item>
        <Form.Item name="password" label="Password" rules={[{ required: true }]}>
          <Input.Password placeholder="Enter your password" />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" block>
            Register
          </Button>
        </Form.Item>
        <Typography.Text>
          Already have an account? <a href="/login">Login</a>
        </Typography.Text>
      </Form>
    </div>
  );
};
