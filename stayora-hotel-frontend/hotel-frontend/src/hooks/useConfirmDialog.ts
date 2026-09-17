import { useCallback, useState } from "react";

export interface ConfirmOptions {
  title: string;
  description: string;
  confirmLabel?: string;
  cancelLabel?: string;
  destructive?: boolean;
}

interface ConfirmState extends ConfirmOptions {
  resolve: (value: boolean) => void;
}

// One reusable confirmation flow for every destructive Admin action
// (delete hotel, delete room, delete image, cancel booking, refund
// payment, delete review) instead of window.confirm() or a bespoke
// dialog per page. Usage: `if (!(await confirm({...}))) return;`
export function useConfirmDialog() {
  const [state, setState] = useState<ConfirmState | null>(null);

  const confirm = useCallback((options: ConfirmOptions) => {
    return new Promise<boolean>((resolve) => {
      setState({ ...options, resolve });
    });
  }, []);

  const handleConfirm = useCallback(() => {
    state?.resolve(true);
    setState(null);
  }, [state]);

  const handleCancel = useCallback(() => {
    state?.resolve(false);
    setState(null);
  }, [state]);

  return { confirm, dialogState: state, handleConfirm, handleCancel };
}
