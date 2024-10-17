import { Wrapper } from "@/components/Wrapper.tsx";
import LogInForm from "@/pages/auth/LogInForm.tsx";

function LogInPage() {
  return (
    <Wrapper>
      <section className="max-w-md mx-auto md:p-4 md:mt-12 md:shadow-lg rounded-lg">
        <h1 className="text-2xl mb-6">Log In</h1>
        <LogInForm />
      </section>
    </Wrapper>
  );
}

export default LogInPage;
