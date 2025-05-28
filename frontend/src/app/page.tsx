'use client'
import { useState } from "react";
import { FiImage, FiUploadCloud } from "react-icons/fi";
import FileShare from "./components/FileShare";

export default function Home() {
  const [activeTab, setActiveTab] = useState<"file" | "image">("file");

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-3xl mx-auto bg-white rounded-xl shadow p-6">
        {/* Tabs */}
        <div className="flex space-x-6 mb-6 ">
          <button
            onClick={() => setActiveTab("file")}
            className={`flex items-center gap-2 pb-3 ${
              activeTab === "file"
                ? "border-b-4 border-[#d10d1a] text-[#d10d1a]"
                : "text-gray-600"
            }`}
          >
            <FiUploadCloud /> File Share
          </button>
          <button
            onClick={() => setActiveTab("image")}
            className={`flex items-center gap-2 pb-3 ${
              activeTab === "image"
                ? "border-b-4 border-[#d10d1a] text-[#d10d1a]"
                : "text-gray-600"
            }`}
          >
            <FiImage /> Image Converter
          </button>
        </div>

        {activeTab === "file" && <FileShare />}
        {activeTab === "image" && (
          <div className="text-center text-gray-500">
            (Image converter coming soon)
          </div>
        )}
      </div>
    </div>
  );
}
