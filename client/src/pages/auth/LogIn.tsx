import { z } from "zod";
import {useAuth} from "../../hooks/useAuth.tsx";
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";

const loginSchema = z.object({
  email: z.string().min(1, "Email is required").email("Invalid email address"),
  password: z.string().min(6, "Password must be at least 6 characters long"),
});

type LoginFormInputs = z.infer<typeof loginSchema>;

function LogIn() {
  const auth = useAuth();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormInputs>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "mikhl90+tania@gmail.com",
      password: "password",
    }
  });

  function onSubmit(credentials: LoginFormInputs) {
    auth.login(credentials)
      .then(() => {
        console.log("logged in!");
      })
      .catch(err => alert(err.message));
  }

  return (
    <section>
      <h1>Log In</h1>
      <form onSubmit={handleSubmit(onSubmit)}>
        <div>
          <label>Email</label>
          <input type="email" {...register("email")} />
          {errors.email && <p>{errors.email.message}</p>}
        </div>

        <div>
          <label>Password</label>
          <input type="password" {...register("password")} />
          {errors.password && <p>{errors.password.message}</p>}
        </div>

        <button type="submit">Login</button>
      </form>
    </section>
  )
}

export default LogIn;
