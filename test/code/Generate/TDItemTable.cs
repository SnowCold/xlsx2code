//Auto Generate Don't Edit it
using UnityEngine;
using System;
using System.IO;
using System.Collections;
using System.Collections.Generic;
using SoulEngine;

namespace LR
{
    public static partial class TDItemTable
    {
        private static TDTableMetaData m_MetaData = new TDTableMetaData(TDItemTable.Parse, "item");
        public static TDTableMetaData metaData
        {
            get { return m_MetaData; }
        }
        
        private static Dictionary<string, TDItem> m_DataCache = new Dictionary<string, TDItem>();
        private static List<TDItem> m_DataList = new List<TDItem >();
        
        public static void Parse(byte[] fileData)
        {
            m_DataCache.Clear();
            m_DataList.Clear();
            DataStreamReader dataR = new DataStreamReader(fileData);
            int rowCount = dataR.GetRowCount();
            int[] fieldIndex = dataR.GetFieldIndex(TDItem.GetFieldHeadIndex());
    #if (UNITY_STANDALONE_WIN) || UNITY_EDITOR || UNITY_STANDALONE_OSX
            dataR.CheckFieldMatch(TDItem.GetFieldHeadIndex(), "ItemTable");
    #endif
            for (int i = 0; i < rowCount; ++i)
            {
                TDItem memberInstance = new TDItem();
                memberInstance.ReadRow(dataR, fieldIndex);
                OnAddRow(memberInstance);
                memberInstance.Reset();
                CompleteRowAdd(memberInstance);
            }
            Debug.Log(string.Format("Parse Success TDItem"));
        }

        private static void OnAddRow(TDItem memberInstance)
        {
            string key = memberInstance.id;
            if (m_DataCache.ContainsKey(key))
            {
                Debug.LogError(string.Format("Invaild,  TDItemTable Id already exists {0}", key));
            }
            else
            {
                m_DataCache.Add(key, memberInstance);
                m_DataList.Add(memberInstance);
            }
        }    
        
        public static void Reload(byte[] fileData)
        {
            Parse(fileData);
        }

        public static int count
        {
            get 
            {
                return m_DataCache.Count;
            }
        }

        public static List<TDItem> dataList
        {
            get 
            {
                return m_DataList;
            }    
        }

        public static TDItem GetData(string key)
        {
            if (m_DataCache.ContainsKey(key))
            {
                return m_DataCache[key];
            }
            else
            {
                Debugger.LogError(string.Format("Can't find key {0} in TDItem", key));
                return null;
            }
        }
    }
}//namespace LR