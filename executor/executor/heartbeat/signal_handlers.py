import sys


class HeartbeatSignalHandlers:
    def signal_term_handler(self, Code, Frame):
        self.logger.info("Received SIGTERM.")
        sys.exit(0)

    def signal_hup_handler(self, signal, frame):
        self.logger.info("Received SIGHUP.")
        sys.exit(0)

    def signal_quit_handler(self, signal, frame):
        self.logger.info("Received SIGQUIT.")
        sys.exit(0)
