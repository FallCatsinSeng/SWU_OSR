"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";
import { useCurrentUser } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useToast } from "@/components/ui/toast";
import { Skeleton } from "@/components/ui/skeleton";

const profileSchema = z.object({
  alias: z.string().min(3, "Alias must be at least 3 characters").max(30),
  bio: z.string().max(500, "Bio must be at most 500 characters").optional(),
  avatar_url: z.string().url("Must be a valid URL").or(z.literal("")).optional(),
});

type ProfileFormValues = z.infer<typeof profileSchema>;

export function ProfileEditForm() {
  const { data: user, isLoading } = useCurrentUser();
  const queryClient = useQueryClient();
  const { toast } = useToast();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ProfileFormValues>({
    resolver: zodResolver(profileSchema),
    values: user
      ? { alias: user.alias, bio: user.bio || "", avatar_url: user.avatar_url || "" }
      : undefined,
  });

  const updateProfile = useMutation({
    mutationFn: async (values: ProfileFormValues) => {
      const { data } = await api.put("/profiles/me", values);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["currentUser"] });
      toast("Profile updated successfully", "success");
    },
    onError: () => {
      toast("Failed to update profile", "error");
    },
  });

  if (isLoading) {
    return <Skeleton className="h-64 w-full" />;
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Edit Profile</CardTitle>
      </CardHeader>
      <CardContent>
        <form
          onSubmit={handleSubmit((data) => updateProfile.mutate(data))}
          className="space-y-4"
        >
          <div className="space-y-2">
            <label htmlFor="alias" className="text-sm font-medium text-gray-700">
              Alias (Public Display Name)
            </label>
            <Input id="alias" {...register("alias")} />
            {errors.alias && (
              <p className="text-sm text-red-600">{errors.alias.message}</p>
            )}
          </div>
          <div className="space-y-2">
            <label htmlFor="bio" className="text-sm font-medium text-gray-700">
              Bio
            </label>
            <Textarea id="bio" {...register("bio")} placeholder="Tell others about yourself..." />
            {errors.bio && (
              <p className="text-sm text-red-600">{errors.bio.message}</p>
            )}
          </div>
          <div className="space-y-2">
            <label htmlFor="avatar_url" className="text-sm font-medium text-gray-700">
              Avatar URL
            </label>
            <Input
              id="avatar_url"
              {...register("avatar_url")}
              placeholder="https://example.com/avatar.png"
            />
            {errors.avatar_url && (
              <p className="text-sm text-red-600">{errors.avatar_url.message}</p>
            )}
          </div>
          <Button type="submit" disabled={updateProfile.isPending}>
            {updateProfile.isPending ? "Saving..." : "Save Changes"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
