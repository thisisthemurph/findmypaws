import { useAuth } from "@clerk/clerk-react";
import { z } from "zod";

interface UseFetchOptions {
  method?: string;
  headers?: HeadersInit;
  body?: BodyInit | null;
}

const errorSchema = z.object({
  message: z.string(),
});

export const useApi = () => {
  const { getToken } = useAuth();

  return async <T = unknown>(url: string, options?: UseFetchOptions): Promise<T extends void ? void : T> => {
    const token = await getToken();

    const headers = {
      Authorization: `Bearer ${token}`,
      ...(options?.body instanceof FormData ? {} : { "Content-Type": "application/json" }),
      ...options?.headers,
    };

    const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}${url}`, {
      ...options,
      headers,
    });

    if (response.status === 204) {
      return undefined as T extends void ? void : T;
    }

    if (!response.ok) {
      const body = await response.json();
      const result = errorSchema.safeParse(body);
      if (result.success) {
        throw new Error(result.data.message);
      }
      throw new Error();
    }

    try {
      const result = await response.json();
      return result as T extends void ? void : T;
    } catch (ex) {
      if (ex instanceof SyntaxError) {
        // Unexpected end of JSON body, no content returned.
        return undefined as T extends void ? void : T;
      }
      throw ex;
    }
  };
};
