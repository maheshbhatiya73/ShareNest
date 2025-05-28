'use client';

import { useState } from 'react';
import { useParams } from 'next/navigation';
import { FiDownloadCloud, FiAlertCircle, FiCheckCircle, FiLock } from 'react-icons/fi';
import { motion } from 'framer-motion';
import { saveAs } from 'file-saver';

const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;

export default function DownloadPage() {
    const { id } = useParams();
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [success, setSuccess] = useState(false);
    const [loading, setLoading] = useState(false);

    const handleDownload = async () => {
        setError('');
        if (!password.trim()) {
            setError('Password is required.');
            return;
        }
        setSuccess(false);
        setLoading(true);
        try {
            console.log('Download started', { id, password, baseUrl });

            const res = await fetch(`${baseUrl}/api/download/${id}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ password }),
            });

            console.log('Response status:', res.status);

            if (!res.ok) {
                const msg = await res.text();
                setError(msg);
                return;
            }

            const blob = await res.blob();
            const contentDisposition = res.headers.get('Content-Disposition');
            const filenameMatch = contentDisposition?.match(/filename="?(.+)"?/);
            const filename = filenameMatch ? filenameMatch[1] : 'file';

            console.log('Filename:', filename);
            saveAs(blob, filename);
            setSuccess(true);
        } catch (err) {
            console.error('Download error:', err);
            setError('Something went wrong.');
        } finally {
            setLoading(false);
        }
    };
    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="max-w-md mx-auto mt-24 bg-gradient-to-br from-white to-slate-100 dark:from-zinc-900 dark:to-zinc-800 p-8 rounded-3xl shadow-2xl"
        >
            <h1 className="text-3xl font-semibold text-gray-800 dark:text-white mb-6 flex items-center gap-3">
                <FiDownloadCloud className="text-indigo-600 dark:text-indigo-400" />
                Secure File Download
            </h1>

            <div className="relative">
                <FiLock className="absolute left-3 top-3 text-gray-400" />
                <input
                    type="password"
                    placeholder="Enter password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full pl-10 pr-4 py-2.5 border rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500 transition text-gray-700 dark:text-white bg-white dark:bg-zinc-900 border-gray-300 dark:border-zinc-700"
                />
            </div>

            <button
                onClick={handleDownload}
                disabled={loading}
                className={`mt-6 w-full flex justify-center items-center gap-2 py-2.5 px-4 rounded-xl font-medium transition 
          ${loading ? 'bg-indigo-400 cursor-not-allowed' : 'bg-indigo-600 hover:bg-indigo-700 cursor-pointer'} 
          text-white shadow-md`}
            >
                {loading ? (
                    <span className="loader w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin" />
                ) : (
                    <>
                        <FiDownloadCloud /> Download File
                    </>
                )}
            </button>

            {error && (
                <div className="mt-4 text-red-600 dark:text-red-400 flex items-center gap-2">
                    <FiAlertCircle /> {error}
                </div>
            )}

            {success && (
                <div className="mt-4 text-green-600 dark:text-green-400 flex items-center gap-2">
                    <FiCheckCircle /> Download started!
                </div>
            )}
        </motion.div>
    );
}
