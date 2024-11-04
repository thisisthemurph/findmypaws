import { Notification } from "@/api/types.ts";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Link } from "react-router-dom";
import { format, formatRelative, isToday, isYesterday, isThisWeek, subWeeks } from "date-fns";
import { useApi } from "@/hooks/useApi.ts";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export default function NotificationMenu({ notifications }: { notifications: Notification[] }) {
  const api = useApi();
  const queryClient = useQueryClient();

  const readAllNotificationsMutation = useMutation({
    mutationFn: () => api("/user/notifications/read-all", { method: "POST" }),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["notifications"] });
    },
    onError: (err: Error) => console.error(err),
  });

  const todayNotifications = notifications.filter((n) => isToday(new Date(n.created_at)));
  const yesterdayNotifications = notifications.filter((n) => isYesterday(new Date(n.created_at)));
  const thisWeekNotifications = notifications.filter(
    (n) =>
      isThisWeek(new Date(n.created_at), { weekStartsOn: 1 }) &&
      !isToday(new Date(n.created_at)) &&
      !isYesterday(new Date(n.created_at))
  );
  const lastWeekNotifications = notifications.filter(
    (n) =>
      new Date(n.created_at) > subWeeks(new Date(), 2) &&
      !isThisWeek(new Date(n.created_at), { weekStartsOn: 1 }) &&
      !isToday(new Date(n.created_at)) &&
      !isYesterday(new Date(n.created_at))
  );
  const olderNotifications = notifications.filter((n) => new Date(n.created_at) <= subWeeks(new Date(), 2));

  const hasUnreadNotifications = notifications.some((n) => !n.seen);

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="secondary" size="icon" className="relative p-0 rounded-full">
          <span>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth={1.5}
              stroke="currentColor"
              className="w-10 h-10"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M14.857 17.082a23.848 23.848 0 0 0 5.454-1.31A8.967 8.967 0 0 1 18 9.75V9A6 6 0 0 0 6 9v.75a8.967 8.967 0 0 1-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 0 1-5.714 0m5.714 0a3 3 0 1 1-5.714 0"
              />
            </svg>
          </span>
          {hasUnreadNotifications && (
            <span className="absolute bottom-0 right-0 w-2.5 h-2.5 bg-red-500 border-2 border-white rounded-full animate-pulse"></span>
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="p-0 w-full mx-2">
        <div className="flex justify-between items-center px-2 py-2 border-b">
          <p className="font-semibold text-slate-700">Notifications</p>
          {hasUnreadNotifications && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => readAllNotificationsMutation.mutate()}
              className="font-normal text-blue-600 hover:text-blue-800"
            >
              Mark all as read
            </Button>
          )}
        </div>
        <div className="p-0 flex flex-col text-sm">
          <NotificationBucket name="Today" notifications={todayNotifications} />
          <NotificationBucket name="Yesterday" notifications={yesterdayNotifications} />
          <NotificationBucket name="This week" notifications={thisWeekNotifications} />
          <NotificationBucket name="Last week" notifications={lastWeekNotifications} />
          <NotificationBucket name="Older" notifications={olderNotifications} />
          {notifications.length === 0 && (
            <div className="flex flex-col gap-4 p-4 text-center">
              <p className="font-semibold text-lg">Nothing to see here</p>
              <p>Come back when there are new notifications.</p>
            </div>
          )}
        </div>
      </PopoverContent>
    </Popover>
  );
}

interface NotificationBucketProps {
  name: string;
  notifications: Notification[];
}

function NotificationBucket({ name, notifications }: NotificationBucketProps) {
  if (notifications.length === 0) return null;

  return (
    <>
      <div className="p-2 tracking-widest font-semibold bg-slate-100 text-slate-700">{name}</div>
      {notifications.map((n) => (
        <Link
          to={n.link}
          key={n.id}
          className="group relative flex gap-4 px-4 py-2 border-b hover:bg-slate-50 transition-colors"
        >
          {!n.seen && <div className="absolute top-2 left-2 w-[6px] h-[6px] bg-blue-600 rounded-full"></div>}
          <div className="flex justify-center items-center bg-slate-100 size-8 rounded-full group-hover:bg-white">
            <MagnifyingGlassIcon className="size-4" />
          </div>
          <div className="flex flex-col gap-2 w-full">
            <p>{n.message}</p>
            <p title={format(n.created_at, "PPP")} className="text-slate-600 text-xs">
              {formatRelative(n.created_at, new Date())}
            </p>
          </div>
        </Link>
      ))}
    </>
  );
}

function MagnifyingGlassIcon({ className }: { className: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      strokeWidth={1.5}
      stroke="currentColor"
      className={className}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"
      />
    </svg>
  );
}
