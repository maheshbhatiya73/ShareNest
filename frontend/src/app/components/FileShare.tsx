import { useState } from "react";
import { motion } from "framer-motion";
import { FiUploadCloud, FiCheckCircle, FiLoader, FiAlertCircle } from "react-icons/fi";

const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
export default function FileShare() {
    const [file, setFile] = useState<File | null>(null);
    const [password, setPassword] = useState("");
    const [uploading, setUploading] = useState(false);
    const [link, setLink] = useState("");
        const [error, setError] = useState("");


    const handleUpload = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(""); // Clear any old errors

        if (!file) {
            setError("Please select a file.");
            return;
        }

        if (!password.trim()) {
            setError("Password is required.");
            return;
        }

        setUploading(true);
        const formData = new FormData();
        formData.append("file", file);
        formData.append("password", password);

        try {
            const res = await fetch(`${baseUrl}/api/upload`, {
                method: "POST",
                body: formData,
            });

            const data = await res.json();
            setLink(data.downloadLink);
        } catch (err) {
            setError("Something went wrong during upload.");
        } finally {
            setUploading(false);
        }
    };

    return (
        <motion.div
            key="file"
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: 20 }}
            transition={{ duration: 0.3 }}
            className="bg-white p-6 rounded-xl shadow-lg max-w-xl mx-auto"
        >
            <h2 className="text-2xl font-bold mb-6 flex items-center gap-2 text-gray-800">
                <FiUploadCloud className="text-[#d10d1a]" /> Share a File
            </h2>

            <form onSubmit={handleUpload} className="space-y-6">
                {/* File Input */}
                <div className="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-[#d10d1a] transition">
                    <input
                        type="file"
                        className="hidden"
                        id="fileInput"
                        onChange={(e) => {
                            const selected = e.target.files?.[0];
                            setFile(selected || null);
                        }}
                    />
                    <label htmlFor="fileInput" className="cursor-pointer text-gray-500 hover:text-[#d10d1a]">
                        {file ? (
                            <span className="font-medium text-gray-800">{file.name}</span>
                        ) : (
                            <span>Click to upload a file</span>
                        )}
                    </label>
                </div>

                {/* Optional Password */}
                <input
                    type="password"
                    placeholder="Password "
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full rounded border px-4 py-2 border-gray-300 focus:outline-none focus:ring-2 focus:ring-[#d10d1a]"
                />

                {/* Upload Button */}
                <button
                    type="submit"
                    disabled={uploading}
                    className="bg-[#d10d1a] text-white w-full py-2 rounded-lg hover:bg-[#d10d1a] flex justify-center items-center gap-2"
                >
                    {uploading ? (
                        <>
                            <FiLoader className="animate-spin" /> Uploading...
                        </>
                    ) : (
                        <>Upload File</>
                    )}
                </button>
                   {error && (
                    <div className="flex items-center gap-2 text-red-600 text-sm">
                        <FiAlertCircle /> {error}
                    </div>
                )}

            </form>

            {/* Success Link */}
            {link && (
                <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    className="mt-6 bg-green-100 text-green-700 p-4 rounded flex flex-col gap-2"
                >
                    <div className="flex items-center gap-2">
                        <FiCheckCircle />
                        <span className="font-medium">Upload successful!</span>
                    </div>
                    <div className="flex items-center justify-between gap-2">
                        <a
                            href={`http://localhost:3000${link}`}
                            className="underline break-all text-sm"
                            target="_blank"
                            rel="noopener noreferrer"
                        >
                            {`http://localhost:3000${link}`}
                        </a>
                        <button
                            className="text-sm text-blue-600 hover:underline"
                            onClick={() => {
                                navigator.clipboard.writeText(`http://localhost:3000${link}`);
                            }}
                        >
                            Copy
                        </button>
                    </div>
                </motion.div>
            )}

        </motion.div>
    );
}
