function Comments() {
    return (
            <div className="chat chat-start">
            {/* Comments inside the CommentsBox.tsx collapsing box*/}
            <div className="chat-image avatar">
                <div className="w-10 rounded-full">
                {/* TODO: Link Profile picture to comment */}
                <img alt="Tailwind CSS chat bubble component" src="https://daisyui.com/images/stock/photo-1534528741775-53994a69daeb.jpg" />
                </div>
            </div>
            <div className="chat-header text-black">
                {/* TODO: Link user name to comment */}
                Placeholder Name
                {/* TODO: Link time to comment*/}
                <time className="text-xs text-black ">12:45</time>
            </div>
                {/* TODO: Link content to comment */}
            <div className="chat-bubble chat-bubble-secondary">Oh no! So terrible!</div>
            </div>  
    );
}

export default Comments;