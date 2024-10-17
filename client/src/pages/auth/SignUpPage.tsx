import { Wrapper } from "@/components/Wrapper";
import SignUpForm from "@/pages/auth/SignUpForm.tsx";

function SignUpPage() {
  return (
    <Wrapper>
      <section className="max-w-md mx-auto md:p-4 md:mt-12 md:shadow-lg rounded-lg">
        <h1 className="text-2xl mb-6">Sign up</h1>
        <SignUpForm />
      </section>
    </Wrapper>
  );
}

export default SignUpPage;
