import React, { useState } from 'react';

function CreatePost() {
    const [selectedFile, setSelectedFile] = useState<File | null>(null);
    const [selectedGroup, setSelectedGroup] = useState<string | null>(null);

    const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        if (event.target.files && event.target.files.length > 0) {
            setSelectedFile(event.target.files[0]);
        }
    };

    const handleGroupSelect = (group: string) => {
        setSelectedGroup(group);
    };

    const handleDeselectGroup = () => {
        setSelectedGroup(null);
    };

    return (
        <div>
            <div className="flex justify-between">
                    <div>
                    {/* Top message box */}
                    <input type="text" placeholder="Title" className="input mt-2 w-full max-w-sm" />
                    </div>
                    
                    <div>
                        <div>
                                {/* Selected Group button*/}
                                {selectedGroup && (
                                <div className="mt-2 btn btn-ghost text-black" onClick={handleDeselectGroup}>
                                    {selectedGroup} X
                                </div>
                                )}

                            <div className="join mt-2">
                                <input className="join-item btn" type="radio" name="options" aria-label="Public" />
                                
                                <input className="join-item btn" type="radio" name="options" aria-label="Private" />
                                
                                <div className="dropdown dropdown-end">
                                    <input className="join-item btn" type="radio" name="options" aria-label="Groups" />
                                    <ul tabIndex={0} className="dropdown-content z-[1] menu p-2 shadow bg-gray-400 rounded-box w-52">
                                        <li>
                                            <a onClick={() => handleGroupSelect('Group1')}>Group1</a>
                                        </li>
                                        <li>
                                            <a onClick={() => handleGroupSelect('Group2')}>Group2</a>
                                        </li>
                                    </ul>
                                
                                </div>
                                
                                <input className="join-item btn" type="radio" name="options" aria-label="Semi-Private" />
                            
                            </div>
                        </div>
                    </div>
            </div>
            {/* Main message box */}
            <div className="relative w-full min-w-[200px] mt-2">
                <textarea
                    className="peer h-full min-h-[200px] w-full resize-none border-b border-blue-gray-200 bg-transparent pt-4 pb-1.5 font-sans text-lg font-normal text-gray-900 outline outline-0 transition-all placeholder-shown:border-blue-gray-200 focus:border-gray-900 focus:outline-0 disabled:resize-none disabled:border-0 disabled:bg-blue-gray-50"
                    placeholder=""
                ></textarea>
                <label
                    className="after:content[' '] pointer-events-none absolute left-0 -top-1.5 flex h-full w-full select-none text-[11px] font-normal leading-tight text-gray-900 transition-all after:absolute after:-bottom-0 after:block after:w-full after:scale-x-0 after:border-b-2 after:border-gray-900 after:transition-transform after:duration-300 peer-placeholder-shown:text-sm peer-placeholder-shown:leading-[4.25] peer-placeholder-shown:text-blue-gray-500 peer-focus:text-[11px] peer-focus:leading-tight peer-focus:text-gray-900 peer-focus:after:scale-x-100 peer-focus:after:border-gray-900 peer-disabled:text-transparent peer-disabled:peer-placeholder-shown:text-blue-gray-500"
                >
                    Message
                </label>
            </div>

            {/* Image upload and preview */}
            <div className="mb-4 flex justify-start items-end">
                <div>
                    <input
                        type="file"
                        id="image-upload"
                        className="hidden"
                        accept="image/*"
                        onChange={handleFileChange}
                    />
                    <label htmlFor="image-upload" className="btn cursor-pointer">
                        Upload Photo
                    </label>
                </div>
                <div className="avatar-preview">
                    {selectedFile && (
                        <img
                            src={URL.createObjectURL(selectedFile)}
                            alt="Preview"
                            className="avatar"
                            style={{ width: 100, height: 100 }}
                        />
                    )}
                </div>

                {/* Post button*/}
                <div className="flex-grow" />
                <button className="btn">Post</button>
            </div>
        </div>
    );
}

export default CreatePost;
