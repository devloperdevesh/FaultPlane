"use client";

import React from "react";
import ErrorState from "./ErrorState";

interface Props {
  children: React.ReactNode;
}

interface State {
  hasError: boolean;
}

export default class ErrorBoundary extends React.Component<Props, State> {
  state = {
    hasError: false,
  };

  static getDerivedStateFromError() {
    return {
      hasError: true,
    };
  }

  componentDidCatch(error: Error) {
    console.error("Dashboard Error:", error);
  }

  render() {
    if (this.state.hasError) {
      return (
        <ErrorState
          retry={() =>
            this.setState({
              hasError: false,
            })
          }
        />
      );
    }

    return this.props.children;
  }
}
