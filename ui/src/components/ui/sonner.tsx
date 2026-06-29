import { Toaster as Sonner, type ToasterProps} from 'sonner'

const toastOptions: ToasterProps['toastOptions'] = {
  classNames: {
    toast: 'cn-toast',
  },
}

const Toaster = ({ ...props }: ToasterProps) => {
  return (
    <Sonner
      theme="light"
      className="toaster group"
      style={
        {
          '--normal-bg': 'var(--popover)',
          '--normal-text': 'var(--popover-foreground)',
          '--normal-border': 'var(--border)',
          '--border-radius': 'var(--radius)',
        } as React.CSSProperties
      }
      toastOptions={toastOptions}
      {...props}
    />
  )
}

export { Toaster }
