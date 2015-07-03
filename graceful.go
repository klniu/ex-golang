package main 

import (

)

// For this we will need to keep track of open connections with a sync.WaitGroup. 
// We will need to increment the wait group on every accepted connection and 
// decrement it on every connection close.
var httpWg sync.WaitGroup

// Here is an example of a listener which increments a wait group on every Accept(). 
// First, we “subclass” net.Listener (you’ll see why we need stop and stopped below):
type gracefulListener struct {
    net.Listener
    stop    chan error
    stopped bool
}

// Next we “override” the Accept method. (Nevermind gracefulConn for now, it will be introduced later).
func (gl *gracefulListener) Accept() (c net.Conn, err error) {
    c, err = gl.Listener.Accept()
    if err != nil {
        return
    }

    c = gracefulConn{Conn: c}

    httpWg.Add(1)
    return
}

// Our Close() method simply sends a nil to the stop channel for the above goroutine to do the rest of the work.
func (gl *gracefulListener) Close() error {
    if gl.stopped {
        return syscall.EINVAL
    }
    gl.stop <- nil
    return <-gl.stop
}

// Finally, this little convenience method extracts the file descriptor from the net.TCPListener.
func (gl *gracefulListener) File() *os.File {
    tl := gl.Listener.(*net.TCPListener)
    fl, _ := tl.File()
    return fl
}


type gracefulConn struct {
    net.Conn
}

// And, of course we also need a variant of a net.Conn which decrements the wait group on Close():
func (w gracefulConn) Close() error {
    httpWg.Done()
    return w.Conn.Close()
}

// The reason the function below starts a goroutine is because this cannot be done in 
// our Accept() above since it will block on gl.Listener.Accept(). The goroutine will 
// unblock it by closing file descriptor.
func newGracefulListener(l net.Listener) (gl *gracefulListener) {
    gl = &gracefulListener{Listener: l, stop: make(chan error)}
    go func() {
        _ = <-gl.stop
        gl.stopped = true
        gl.stop <- gl.Listener.Close()
    }()
    return
}


func main() {
	// And there is one more thing. You should avoid hanging connections that the client 
	// has no intention of closing (or not this week). It is better to create your server as follows:
	server := &http.Server{
        Addr:           "0.0.0.0:8888",
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 16}



	// To start using the above graceful version of the Listener, all we need is to change the server.Serve(l) line to:
	netListener = newGracefulListener(l)
	server.Serve(netListener)

}